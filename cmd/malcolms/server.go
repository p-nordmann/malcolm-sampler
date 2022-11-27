package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	UUID "github.com/google/uuid"

	ma "github.com/p-nordmann/malcolm-sampler"
	pb "github.com/p-nordmann/malcolm-sampler/grpc"
)

// TODO(p-nordmann): do not embed internal errors in grpc responses.

type store struct {
	boundaries map[string]ma.Boundaries
	factories  map[string]ma.SamplerFactory
	lock       sync.Mutex
}

func (s *store) setB(uuid string, b ma.Boundaries) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.boundaries[uuid] = b
}

func (s *store) setF(uuid string, f ma.SamplerFactory) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.factories[uuid] = f
}

func (s *store) getB(uuid string) (ma.Boundaries, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	b, ok := s.boundaries[uuid]
	return b, ok
}

func (s *store) getF(uuid string) (ma.SamplerFactory, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	f, ok := s.factories[uuid]
	return f, ok
}

type samplingServer struct {
	pb.UnimplementedMalcolmSamplerServer

	state store
}

func NewServer() *samplingServer {
	return &samplingServer{
		state: store{
			boundaries: make(map[string]ma.Boundaries),
			factories:  make(map[string]ma.SamplerFactory),
		},
	}
}

func (s *samplingServer) AddBoundaries(ctx context.Context, boundariesMessage *pb.Boundaries) (*pb.BoundariesUUID, error) {
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
	uuid := UUID.New().String()
	s.state.setB(uuid, ma.Boundaries{
		Infima:  boundariesMessage.Infima,
		Suprema: boundariesMessage.Suprema,
	})
	return &pb.BoundariesUUID{Value: uuid}, nil
}

func (s *samplingServer) AddPosterior(sampleStream pb.MalcolmSampler_AddPosteriorServer) error {

	var uuid string
	var boundaries ma.Boundaries
	var samples [][]float64
	var posteriorValues []float64

	for {
		// Receive message and validate inputs.
		msg, err := sampleStream.Recv()

		// When stream terminates, finish building factory.
		if err == io.EOF {
			if len(uuid) == 0 {
				return errors.New("empty query")
			}
			if len(samples) == 0 {
				return errors.New("no data provided")
			}
			factoryUuid := UUID.New().String()
			factory, err := ma.FromSamples(boundaries, samples, posteriorValues)
			if err != nil {
				return fmt.Errorf("error creating factory: %w", err)
			}
			s.state.setF(factoryUuid, factory)
			return sampleStream.SendAndClose(&pb.PosteriorUUID{Value: factoryUuid})
		}

		// Process error cases.
		if err != nil {
			return fmt.Errorf("error receiving true samples: %w", err)
		}
		msgUuid := msg.GetUuid().GetValue()
		if len(msgUuid) == 0 {
			return errors.New("empty UUID is not allowed")
		}

		// First message defines what boundaries to work with.
		if len(uuid) == 0 {
			var ok bool
			uuid = msgUuid
			if boundaries, ok = s.state.getB(uuid); !ok {
				return errors.New("invalid UUID: no corresponding boundaries")
			}
		}

		// Validate UUID and samples dimensions.
		dimension := len(boundaries.Infima)
		if uuid != msgUuid {
			return fmt.Errorf("inconsistent UUID: expected %s but received %s", uuid, msgUuid)
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

func (s *samplingServer) MakeSamples(msg *pb.MakeSamplesRequest, sampleStream pb.MalcolmSampler_MakeSamplesServer) error {
	uuid := msg.GetUuid().GetValue()
	amount := int(msg.GetAmount())
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if factory, ok := s.state.getF(uuid); ok {
		sampler, err := factory.NewSampler(msg.GetOrigin())
		if err != nil {
			return fmt.Errorf("error building sampler: %w", err)
		}
		for k := 0; k < amount; k++ {
			sampleStream.Send(&pb.SamplesBatch{Coordinates: sampler.Sample()})
		}
		return nil
	}
	return errors.New("invalid UUID")
}
