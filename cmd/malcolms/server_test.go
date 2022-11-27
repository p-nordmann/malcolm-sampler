package main_test

import (
	"context"
	"fmt"
	"io"
	"sync"
	"testing"

	m "github.com/p-nordmann/malcolm-sampler"
	malcolms "github.com/p-nordmann/malcolm-sampler/cmd/malcolms"
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

// mockSamplesStreamChan provides a similar interface to `mockSamplesStream` but returns samples
// through a channel.
type mockSamplesStreamChan struct {
	dimension int
	samples   chan []float64
	grpc.ServerStream
}

func makeMockSamplesStreamChan(dimension int) *mockSamplesStreamChan {
	return &mockSamplesStreamChan{
		dimension: dimension,
		samples:   make(chan []float64),
	}
}

func (s *mockSamplesStreamChan) Send(m *mgrpc.SamplesBatch) error {
	for k := 0; k < len(m.Coordinates); k += s.dimension {
		s.samples <- m.Coordinates[k : k+s.dimension]
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

// server will be used by the tests.
var server = malcolms.NewServer()

// sendBoundaries builds a grpc message and calls AddBoundaries.
//
// It returns the value of the UUID and the error returned by the server.
func sendBoundaries(b m.Boundaries) (string, error) {
	dimension := len(b.Infima)
	boundariesUUID, err := server.AddBoundaries(
		context.Background(),
		&mgrpc.Boundaries{
			Dimension: int64(dimension),
			Infima:    b.Infima,
			Suprema:   b.Suprema,
		},
	)
	if boundariesUUID == nil {
		return "", err
	}
	return boundariesUUID.Value, err
}

// We split the logic of  sendPosterior in multiple functions to be able to keep multiple
// streams opened at the same time.
type sendPosteriorInternals struct {
	// Function args.
	f              flatPosterior
	batcheSizes    [][2]int
	boundariesUUID string

	// Local variable.
	stream *mockAddPosteriorStream
}

// Launches AddPosterior stream.
func (sp *sendPosteriorInternals) launch() {
	sp.stream = makeMockAddPosteriorStream()
	go func() { sp.stream.err <- server.AddPosterior(sp.stream) }()
}

// Send posterior in batches.
func (sp *sendPosteriorInternals) sendBatches() {
	coordinates := sp.f.coordinates
	posteriorValues := sp.f.posteriorValues
	for _, batchSize := range sp.batcheSizes {
		sp.stream.next <- &mgrpc.PosteriorValuesBatch{
			Uuid:            &mgrpc.BoundariesUUID{Value: sp.boundariesUUID},
			Coordinates:     coordinates[:batchSize[0]],
			PosteriorValues: posteriorValues[:batchSize[1]],
		}
		coordinates = coordinates[batchSize[0]:]
		posteriorValues = posteriorValues[batchSize[1]:]
	}
}

// Closes, waits for completion and returns.
func (sp *sendPosteriorInternals) finish() (string, error) {
	close(sp.stream.next)
	select {
	case err := <-sp.stream.err:
		return "", err
	case result := <-sp.stream.result:
		return result.Value, nil
	}
}

// sendPosterior builds a grpc stream and calls AddPosterior.
//
// It returns the value of the UUID and the error returned by the server.
func sendPosterior(f flatPosterior, batcheSizes [][2]int, boundariesUUID string) (string, error) {
	sp := &sendPosteriorInternals{
		f:              f,
		batcheSizes:    batcheSizes,
		boundariesUUID: boundariesUUID,
	}
	sp.launch()
	sp.sendBatches()
	return sp.finish()
}

// The server should allow to add boundaries, posterior and make samples.
//
// We test that it can be done with posterior values sent in batches of 1, 2 or 3.
//
// We also validate that it can be achieved in parallel with two rpc calls.
//
// TODO: mock sampler to make sure boundaries are correctly sliced back to 2D.
func TestBasicUseCase(t *testing.T) {

	// Example values for tests.
	exampleBoundaries := m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 1}}
	exampleFlatPosterior := toRowMajor(posterior{
		coordinates: [][]float64{
			{0.1, 0.1, 0.9},
			{0.9, 0.1, 0.1},
			{0.1, 0.1, 0.9},
			{0.1, 0.1, 0.1},
			{0.1, 0.1, 0.9},
			{0.1, 0.1, 0.1}},
		posteriorValues: []float64{1, 2, 3, 4, 5, 6},
	})
	exampleOrigin := []float64{0.5, 0.5, 0.5}

	// Fixture data.
	basicTestCases := []struct {
		dimension  int
		boundaries m.Boundaries
		posterior  flatPosterior
		batchSizes [][2]int
		origin     []float64
		amount     int
	}{
		{
			dimension:  3,
			boundaries: exampleBoundaries,
			posterior:  exampleFlatPosterior,
			batchSizes: [][2]int{{3, 1}, {3, 1}, {3, 1}, {3, 1}, {3, 1}, {3, 1}},
			origin:     exampleOrigin,
			amount:     10,
		},
		{
			dimension:  3,
			boundaries: exampleBoundaries,
			posterior:  exampleFlatPosterior,
			batchSizes: [][2]int{{6, 2}, {6, 2}, {6, 2}},
			origin:     exampleOrigin,
			amount:     10,
		},
		{
			dimension:  3,
			boundaries: exampleBoundaries,
			posterior:  exampleFlatPosterior,
			batchSizes: [][2]int{{9, 3}, {9, 3}},
			origin:     exampleOrigin,
			amount:     10,
		},
	}

	// Run test cases.
	for k, tc := range basicTestCases {
		t.Run(
			fmt.Sprintf(
				"#%d: should allow to perform typical user story",
				k,
			),
			func(t *testing.T) {
				boundariesUUID, _ := sendBoundaries(tc.boundaries)
				posteriorUUID, _ := sendPosterior(tc.posterior, tc.batchSizes, boundariesUUID)

				stream := makeMockSamplesStream(tc.dimension)
				server.MakeSamples(
					&mgrpc.MakeSamplesRequest{
						Uuid:   &mgrpc.PosteriorUUID{Value: posteriorUUID},
						Origin: tc.origin,
						Amount: int64(tc.amount),
					},
					stream,
				)

				if len(stream.samples) != tc.amount {
					t.Errorf("expected %d samples, got %d", tc.amount, len(stream.samples))
				}
			},
		)
	}
}

// The service should be able to handle parallel rpc calls.
//
// Test this with -race option.
func TestParallelCalls(t *testing.T) {
	t.Run("should handle parallel calls to AddPosterior", func(t *testing.T) {
		boundariesUUID, _ := sendBoundaries(m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 1}})

		sp1 := &sendPosteriorInternals{
			f:              flatPosterior{coordinates: []float64{0.5, 0.5, 0.5}, posteriorValues: []float64{1}},
			batcheSizes:    [][2]int{{3, 1}},
			boundariesUUID: boundariesUUID,
		}
		sp2 := &sendPosteriorInternals{
			f:              flatPosterior{coordinates: []float64{0.5, 0.5, 0.5, 1, 1, 1}, posteriorValues: []float64{1, 2}},
			batcheSizes:    [][2]int{{3, 1}, {3, 1}},
			boundariesUUID: boundariesUUID,
		}

		var wg sync.WaitGroup
		var u1, u2 string
		var err1, err2 error

		wg.Add(2)
		go func() {
			sp1.launch()
			sp1.sendBatches()
			wg.Done()
		}()
		go func() {
			sp2.launch()
			sp2.sendBatches()
			wg.Done()
		}()
		wg.Wait()

		wg.Add(2)
		go func() {
			u1, err1 = sp1.finish()
			wg.Done()
		}()
		go func() {
			u2, err2 = sp2.finish()
			wg.Done()
		}()
		wg.Wait()

		if err1 != nil || err2 != nil {
			t.Errorf("expected successful calls to AddPosterior but got errors %v, %v", err1, err2)
		}
		if u1 == u2 {
			t.Errorf("both UUIDs were equal ('%s')", u1)
		}
	})

	t.Run("should handle parallel calls to MakeSamples", func(t *testing.T) {
		boundariesUUID, _ := sendBoundaries(m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 1}})
		exampleFlatPosterior := flatPosterior{coordinates: []float64{0.5, 0.5, 0.5}, posteriorValues: []float64{1}}
		posteriorUUID, _ := sendPosterior(exampleFlatPosterior, [][2]int{{3, 1}}, boundariesUUID)

		var wg sync.WaitGroup
		wg.Add(2)
		var err1, err2 error

		s1 := makeMockSamplesStreamChan(3)
		go func() {
			err1 = server.MakeSamples(
				&mgrpc.MakeSamplesRequest{
					Uuid:   &mgrpc.PosteriorUUID{Value: posteriorUUID},
					Origin: []float64{0.9, 0.9, 0.9},
					Amount: 10,
				},
				s1,
			)
			wg.Done()
		}()

		s2 := makeMockSamplesStreamChan(3)
		go func() {
			err2 = server.MakeSamples(
				&mgrpc.MakeSamplesRequest{
					Uuid:   &mgrpc.PosteriorUUID{Value: posteriorUUID},
					Origin: []float64{0.1, 0.1, 0.1},
					Amount: 10,
				},
				s2,
			)
			wg.Done()
		}()

		// Receive samples.
		for k := 0; k < 10; k++ {
			<-s1.samples
			<-s2.samples
		}
		wg.Wait()
		if err1 != nil || err2 != nil {
			t.Errorf("expected successful calls to MakeSamples but got errors %v, %v", err1, err2)
		}
	})
}

// Multiple calls to AddBoundaries or AddPosterior should create different UUIDs.
//
// Of course, this test is very shallow but should avoid the basic mistake of hardcoding a UUID value.
func TestUUIDAreUnique(t *testing.T) {
	exampleBoundaries := m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 1}}
	exampleFlatPosterior := flatPosterior{coordinates: []float64{0.5, 0.5, 0.5}, posteriorValues: []float64{1}}

	t.Run("two calls to AddBoundaries should return different uuids", func(t *testing.T) {
		u1, _ := sendBoundaries(exampleBoundaries)
		u2, _ := sendBoundaries(exampleBoundaries)
		if u1 == u2 {
			t.Errorf("both UUIDs were equal ('%s')", u1)
		}
	})

	t.Run("two calls to AddPosterior should return different uuids", func(t *testing.T) {
		boundariesUUID, _ := sendBoundaries(exampleBoundaries)
		u1, _ := sendPosterior(exampleFlatPosterior, [][2]int{{3, 1}}, boundariesUUID)
		u2, _ := sendPosterior(exampleFlatPosterior, [][2]int{{3, 1}}, boundariesUUID)
		if u1 == u2 {
			t.Errorf("both UUIDs were equal ('%s')", u1)
		}
	})
}

// Basic failure cases should be gracefully handled and trigger nice errors from the server.
func TestFailureCases(t *testing.T) {
	t.Run("posterior should be expected in row-major order", func(t *testing.T) {
		boundariesUUID, _ := sendBoundaries(m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 3}})
		exampleFlatPosterior := toColumnMajor(
			posterior{
				coordinates: [][]float64{
					{0.5, 0.5, 2.5},
					{0.5, 0.5, 2.5},
					{0.5, 0.5, 2.5},
					{0.5, 0.5, 2.5},
				},
				posteriorValues: []float64{1, 2, 3, 4},
			},
		)
		_, err := sendPosterior(exampleFlatPosterior, [][2]int{{12, 4}}, boundariesUUID)
		if err == nil {
			t.Error("Expected error out of bounds but got <nil>.")
		}
	})
	t.Run("should fail when providing wrong UUID to AddPosterior", func(t *testing.T) {
		boundariesUUID, _ := sendBoundaries(m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 1}})
		exampleFlatPosterior := flatPosterior{coordinates: []float64{0.5, 0.5, 0.5}, posteriorValues: []float64{1}}
		_, err := sendPosterior(exampleFlatPosterior, [][2]int{{3, 1}}, boundariesUUID+"-wrong-uuid")
		if err == nil {
			t.Error("expected error 'invalid UUID' but got <nil>")
		}
	})
	t.Run("should fail when providing wrong UUID to MakeSamples", func(t *testing.T) {
		boundariesUUID, _ := sendBoundaries(m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 1}})
		exampleFlatPosterior := flatPosterior{coordinates: []float64{0.5, 0.5, 0.5}, posteriorValues: []float64{1}}
		posteriorUUID, _ := sendPosterior(exampleFlatPosterior, [][2]int{{3, 1}}, boundariesUUID)

		stream := makeMockSamplesStream(3)
		err := server.MakeSamples(
			&mgrpc.MakeSamplesRequest{
				Uuid:   &mgrpc.PosteriorUUID{Value: posteriorUUID + "-wrong-uuid"},
				Origin: []float64{0.5, 0.5, 0.5},
				Amount: 10,
			},
			stream,
		)

		if err == nil {
			t.Error("expected error 'invalid UUID' but got <nil>")
		}
	})
	t.Run("should fail when providing posterior coordinate out of bounds", func(t *testing.T) {
		boundariesUUID, _ := sendBoundaries(m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 1}})
		outOfBoundsFlatPosterior := flatPosterior{coordinates: []float64{0.5, 1.5, 0.5}, posteriorValues: []float64{1}}
		_, err := sendPosterior(outOfBoundsFlatPosterior, [][2]int{{3, 1}}, boundariesUUID)
		if err == nil {
			t.Error("expected error 'out-of-bounds' but got <nil>")
		}
	})
	t.Run("should fail when providing posterior coordinates of incorrect dimension", func(t *testing.T) {
		boundariesUUID, _ := sendBoundaries(m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 1}})
		exampleFlatPosterior := flatPosterior{coordinates: []float64{0.5, 0.5, 0.5, 1, 1, 1, 1}, posteriorValues: []float64{1, 2}}
		_, err := sendPosterior(exampleFlatPosterior, [][2]int{{3, 1}, {4, 1}}, boundariesUUID)
		if err == nil {
			t.Error("expected error 'invalid dimension' but got <nil>")
		}
	})
}
