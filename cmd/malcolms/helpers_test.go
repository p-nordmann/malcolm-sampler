package main_test

import (
	"io"

	mgrpc "github.com/p-nordmann/malcolm-sampler/grpc"
	"google.golang.org/grpc"
)

// mockAddPosteriorStream provides mocking for request's stream for AddPosterior.
type mockAddPosteriorStream struct {
	next   chan *mgrpc.PosteriorValuesBatch
	result chan *mgrpc.PosteriorUUID
	err    chan error
	grpc.ServerStream
}

func makeMockAddPosteriorStream() *mockAddPosteriorStream {
	return &mockAddPosteriorStream{
		next:   make(chan *mgrpc.PosteriorValuesBatch),
		result: make(chan *mgrpc.PosteriorUUID),
		err:    make(chan error),
	}
}

func (s *mockAddPosteriorStream) SendAndClose(uuid *mgrpc.PosteriorUUID) error {
	s.result <- uuid
	close(s.result)
	return nil
}

func (s *mockAddPosteriorStream) Recv() (*mgrpc.PosteriorValuesBatch, error) {
	next, ok := <-s.next
	if !ok {
		return nil, io.EOF
	}
	return next, nil
}

// mockSamplesStream provides mocking for response's stream for MakeSamples.
//
// TODO: add the possibility to stall on Send in order to synchronize parallel calls to MakeSamples.
type mockSamplesStream struct {
	dimension int
	samples   [][]float64
	grpc.ServerStream
}

func makeMockSamplesStream(dimension int) *mockSamplesStream {
	return &mockSamplesStream{
		dimension: dimension,
	}
}

func (s *mockSamplesStream) Send(m *mgrpc.SamplesBatch) error {
	for k := 0; k < len(m.Coordinates); k += s.dimension {
		s.samples = append(s.samples, m.Coordinates[k:k+s.dimension])
	}
	return nil
}

// posterior describes a set of points with posterior values.
type posterior struct {
	coordinates     [][]float64
	posteriorValues []float64
}

// flatPosterior represents the same information as posterior struct with flattened coordinates.
type flatPosterior struct {
	coordinates     []float64
	posteriorValues []float64
}

// toRowMajor converts a posterior to flatPosterior by flattening coordinates in row-major order.
func toRowMajor(p posterior) flatPosterior {
	f := flatPosterior{
		posteriorValues: p.posteriorValues,
	}
	if len(p.coordinates) > 0 {
		dimension := len(p.coordinates[0])
		f.coordinates = make([]float64, len(p.coordinates)*dimension)
		for k, point := range p.coordinates {
			for i, coordinate := range point {
				f.coordinates[k*dimension+i] = coordinate
			}
		}
	}
	return f
}

// toColumnMajor converts a posterior to flatPosterior by flattening coordinates in column-major order.
func toColumnMajor(p posterior) flatPosterior {
	f := flatPosterior{
		posteriorValues: p.posteriorValues,
	}
	count := len(p.coordinates)
	if count > 0 {
		f.coordinates = make([]float64, count*len(p.coordinates[0]))
		for k, point := range p.coordinates {
			for i, coordinate := range point {
				f.coordinates[i*count+k] = coordinate
			}
		}
	}
	return f
}
