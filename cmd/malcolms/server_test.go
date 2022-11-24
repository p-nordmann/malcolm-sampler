package main_test

import (
	"context"
	"testing"

	m "github.com/p-nordmann/malcolm-sampler"
	malcolms "github.com/p-nordmann/malcolm-sampler/cmd/malcolms"
	mgrpc "github.com/p-nordmann/malcolm-sampler/grpc"
)

// server will be used by the tests.
var server = malcolms.NewServer()

// sendBoundaries builds a grpc message and triggers AddBoundaries.
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

// sendPosterior builds a grpc stream and triggers AddPosterior.
//
// It returns the value of the UUID and the error returned by the server.
func sendPosterior(f flatPosterior, dimension, batchSize int, boundariesUUID string) (string, error) {
	stream := makeMockAddPosteriorStream()
	go func() { stream.err <- server.AddPosterior(stream) }()
	for batchCount := 0; batchCount*batchSize < len(f.posteriorValues); batchCount++ {
		stream.next <- &mgrpc.PosteriorValuesBatch{
			Uuid:            &mgrpc.BoundariesUUID{Value: boundariesUUID},
			Coordinates:     f.coordinates[batchCount*dimension*batchSize : (batchCount+1)*dimension*batchSize],
			PosteriorValues: f.posteriorValues[batchCount*batchSize : (batchCount+1)*batchSize],
		}
	}
	close(stream.next)
	select {
	case err := <-stream.err:
		return "", err
	case result := <-stream.result:
		return result.Value, nil
	}
}

// The server should allow to add boundaries, posterior and make samples.
//
// We test that it can be done with posterior values sent in batches of 1, 2 or 3.
//
// We also validate that it can be achieved in parallel with two rpc calls.
//
// TODO: mock sampler to make sure boundaries are correctly sliced back to 2D.
func TestBasicUseCase(t *testing.T) {
	dimension := 3
	exampleBoundaries := m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 1}}

	t.Run("should allow to register boundaries, posterior (batch_size=1) and make samples", func(t *testing.T) {
		exampleFlatPosterior := toRowMajor(posterior{
			coordinates:     [][]float64{{0.1, 0.1, 0.9}, {0.9, 0.1, 0.1}, {0.1, 0.1, 0.9}},
			posteriorValues: []float64{1, 2, 3},
		})

		boundariesUUID, _ := sendBoundaries(exampleBoundaries)
		posteriorUUID, _ := sendPosterior(exampleFlatPosterior, dimension, 1, boundariesUUID)
		stream := makeMockSamplesStream(dimension)
		server.MakeSamples(
			&mgrpc.MakeSamplesRequest{
				Uuid:   &mgrpc.PosteriorUUID{Value: posteriorUUID},
				Origin: []float64{0.5, 0.5, 0.5},
				Amount: 10,
			},
			stream,
		)

		if len(stream.samples) != 10 {
			t.Errorf("expected 10 samples, got %d", len(stream.samples))
		}
	})

	t.Run("should allow to register boundaries, posterior (batch_size=2) and make samples", func(t *testing.T) {
		exampleFlatPosterior := toRowMajor(posterior{
			coordinates:     [][]float64{{0.1, 0.1, 0.9}, {0.9, 0.1, 0.1}, {0.1, 0.1, 0.9}, {0.1, 0.1, 0.1}},
			posteriorValues: []float64{1, 2, 3, 4},
		})

		boundariesUUID, _ := sendBoundaries(exampleBoundaries)
		posteriorUUID, _ := sendPosterior(exampleFlatPosterior, dimension, 2, boundariesUUID)
		stream := makeMockSamplesStream(dimension)
		server.MakeSamples(
			&mgrpc.MakeSamplesRequest{
				Uuid:   &mgrpc.PosteriorUUID{Value: posteriorUUID},
				Origin: []float64{0.5, 0.5, 0.5},
				Amount: 10,
			},
			stream,
		)

		if len(stream.samples) != 10 {
			t.Errorf("expected 10 samples, got %d", len(stream.samples))
		}
	})

	t.Run("should allow to register boundaries, posterior (batch_size=3) and make samples", func(t *testing.T) {
		exampleFlatPosterior := toRowMajor(posterior{
			coordinates:     [][]float64{{0.1, 0.1, 0.9}, {0.9, 0.1, 0.1}, {0.1, 0.1, 0.9}},
			posteriorValues: []float64{1, 2, 3},
		})

		boundariesUUID, _ := sendBoundaries(exampleBoundaries)
		posteriorUUID, _ := sendPosterior(exampleFlatPosterior, dimension, 3, boundariesUUID)
		stream := makeMockSamplesStream(dimension)
		server.MakeSamples(
			&mgrpc.MakeSamplesRequest{
				Uuid:   &mgrpc.PosteriorUUID{Value: posteriorUUID},
				Origin: []float64{0.5, 0.5, 0.5},
				Amount: 10,
			},
			stream,
		)

		if len(stream.samples) != 10 {
			t.Errorf("expected 10 samples, got %d", len(stream.samples))
		}
	})

	t.Run("should allow to handle parallel calls to makeSamples", func(t *testing.T) {

	})
}

// The service should be able to handle parallel calls to AddPosterior.
func TestParallelAddPosterior(t *testing.T) {
	t.Run("should allow to register several posteriors in parallel", func(t *testing.T) {

	})
}

// Multiple calls to AddBoundaries or AddPosterior should create different UUIDs.
//
// Of course, this test is very shallow but should avoid the basic mistake of hardcoding a UUID value.
func TestUUIDAreUnique(t *testing.T) {
	dimension := 3
	exampleBoundaries := m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 1}}
	exampleFlatPosterior := toRowMajor(posterior{coordinates: [][]float64{{0.5, 0.5, 0.5}}, posteriorValues: []float64{1}})

	t.Run("two calls to AddBoundaries should return different uuids", func(t *testing.T) {
		u1, _ := sendBoundaries(exampleBoundaries)
		u2, _ := sendBoundaries(exampleBoundaries)
		if u1 == u2 {
			t.Errorf("both UUIDs were equal ('%s')", u1)
		}
	})

	t.Run("two calls to AddPosterior should return different uuids", func(t *testing.T) {
		boundariesUUID, _ := sendBoundaries(exampleBoundaries)
		u1, _ := sendPosterior(exampleFlatPosterior, dimension, 1, boundariesUUID)
		u2, _ := sendPosterior(exampleFlatPosterior, dimension, 1, boundariesUUID)
		if u1 == u2 {
			t.Errorf("both UUIDs were equal ('%s')", u1)
		}
	})
}

// Basic failure cases should be gracefully handled and trigger nice errors from the server.
func TestFailureCases(t *testing.T) {
	t.Run("posterior should be expected in row-major order", func(t *testing.T) {
		dimension := 3
		exampleBoundaries := m.Boundaries{Infima: []float64{0, 0, 0}, Suprema: []float64{1, 1, 3}}
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
		boundariesUUID, _ := sendBoundaries(exampleBoundaries)
		_, err := sendPosterior(exampleFlatPosterior, dimension, 4, boundariesUUID)
		if err == nil {
			t.Error("Expected error out of bounds but got <nil>.")
		}
	})
	t.Run("should fail when providing wrong UUID to AddPosterior", func(t *testing.T) {

	})
	t.Run("should fail when providing wrong UUID to MakeSamples", func(t *testing.T) {

	})
	t.Run("should fail when providing posterior coordinate out of bounds", func(t *testing.T) {

	})
	t.Run("should fail when providing posterior coordinates of incorrect dimension", func(t *testing.T) {

	})
}
