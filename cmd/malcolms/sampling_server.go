package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"

	m "github.com/p-nordmann/malcolm-sampler"
	pb "github.com/p-nordmann/malcolm-sampler/grpc"
)

// TODO(p-nordmann): do not embed internal errors in grpc responses.

// TODO(p-nordmann): thread safety.
type store struct {
	boundaries map[string]m.Boundaries
	factories  map[string]m.SamplerFactory
}

type samplingServer struct {
	pb.UnimplementedMalcolmSamplerServer

	state store
}

func (s *samplingServer) PutBoundaries(ctx context.Context, boundariesMessage *pb.Boundaries) (*pb.BoundariesUUID, error) {
	// Validate message.
	dimension := int(boundariesMessage.Dimension)
	if boundariesMessage.Dimension < 1 {
		return nil, errors.New("dimension must be positive")
	}
	if len(boundariesMessage.Infima) != dimension || len(boundariesMessage.Suprema) != dimension {
		return nil, errors.New("infima and suprema must have correct dimension")
	}
	for k := 0; k < dimension; k++ {
		if boundariesMessage.Suprema[k] <= boundariesMessage.Infima[k] {
			return nil, errors.New("upper bounds must be greater than lower bounds")
		}
	}
	// Build and store boundaries.
	UUID := uuid.New().String()
	s.state.boundaries[UUID] = m.Boundaries{
		Infima:  boundariesMessage.Infima,
		Suprema: boundariesMessage.Suprema,
	}
	return &pb.BoundariesUUID{Value: UUID}, nil
}

func (s *samplingServer) AddPosterior(sampleStream pb.MalcolmSampler_AddPosteriorServer) error {

	var UUID string
	var boundaries m.Boundaries
	var samples [][]float64
	var posteriorValues []float64

	for {
		// Receive message and validate inputs.
		msg, err := sampleStream.Recv()

		// When stream terminates, finish building factory.
		if err == io.EOF {
			if len(UUID) == 0 {
				return errors.New("empty query")
			}
			if len(samples) == 0 {
				return errors.New("no data provided")
			}
			factoryUUID := uuid.New().String()
			factory, err := m.FromSamples(boundaries, samples, posteriorValues)
			if err != nil {
				return fmt.Errorf("error creating factory: %w", err)
			}
			s.state.factories[factoryUUID] = factory
			return sampleStream.SendAndClose(&pb.PosteriorUUID{Value: factoryUUID})
		}

		// Process error cases.
		if err != nil {
			return fmt.Errorf("error receiving true samples: %w", err)
		}
		msgUUID := msg.GetUuid().GetValue()
		if len(msgUUID) == 0 {
			return errors.New("empty UUID is not allowed")
		}

		// First message defines what boundaries to work with.
		if len(UUID) == 0 {
			var ok bool
			UUID = msgUUID
			if boundaries, ok = s.state.boundaries[UUID]; !ok {
				return errors.New("invalid UUID: no corresponding boundaries")
			}
		}

		// Validate UUID and samples dimensions.
		dimension := len(boundaries.Infima)
		if UUID != msgUUID {
			return fmt.Errorf("inconsistent UUID: expected %s but received %s", UUID, msgUUID)
		}
		coordinates := msg.GetCoordinates()
		posterior := msg.GetPosteriorValues()
		if len(coordinates) == 0 || len(posterior) == 0 {
			return errors.New("messages must provide a positive number of samples")
		}
		if len(coordinates)%dimension != 0 {
			return errors.New("len(coordinates) must be a multiple of space's dimension")
		}
		if len(posterior)*dimension != len(coordinates) {
			return errors.New("incorrect number of posterior values")
		}

		// Validate boundaries and concatenate.
		for k := range posterior {
			point := coordinates[k*dimension : (k+1)*dimension]
			// Validate that samples are not out of bounds.
			for i := range point {
				if point[i] < boundaries.Infima[i] || boundaries.Suprema[i] < point[i] {
					return errors.New("sample out of bounds")
				}
			}
			samples = append(samples, point)
			posteriorValues = append(posteriorValues, posterior[k])
		}
	}
}

func (s *samplingServer) Walk(msg *pb.MakeSamplesRequest, sampleStream pb.MalcolmSampler_MakeSamplesServer) error {
	UUID := msg.GetUuid().GetValue()
	numberOfSamples := int(msg.GetAmount())
	if numberOfSamples <= 0 {
		return errors.New("number_of_samples must be positive")
	}
	if factory, ok := s.state.factories[UUID]; ok {
		sampler, err := factory.NewSampler(msg.GetOrigin())
		if err != nil {
			return fmt.Errorf("error building sampler: %w", err)
		}
		for k := 0; k < numberOfSamples; k++ {
			sampleStream.Send(&pb.SamplesBatch{Coordinates: sampler.Sample()})
		}
		return nil
	}
	return errors.New("invalid UUID")
}
