package malcolms

import (
	"errors"
)

// FromSamples runs preprocessing on samples and posterior values and returns a structure that
// implements TruePosterior interface.
//
// Dimensions are expected to be consistent between all parameters.
func FromSamples(
	boundaries Boundaries, samples [][]float64, posteriorValues []float64,
) (SamplerFactory, error) {
	dimension := boundaries.dimension()

	// Validate boundaries.
	if len(boundaries.Infima) != len(boundaries.Suprema) || dimension == 0 {
		return nil, errors.New("invalid boundaries")
	}
	for i := 0; i < dimension; i++ {
		if boundaries.Suprema[i] <= boundaries.Infima[i] {
			return nil, errors.New("malformed boundaries")
		}
	}

	// Validate samples.
	if len(samples) == 0 {
		return nil, errors.New("no sample")
	}
	for _, sample := range samples {
		if len(sample) != dimension {
			return nil, errors.New("sample with invalid dimension")
		}
		for i, coordinate := range sample {
			if coordinate < boundaries.Infima[i] || boundaries.Suprema[i] < coordinate {
				return nil, errors.New("sample out of bounds")
			}
		}
	}

	// Validate posterior values.
	if len(samples) != len(posteriorValues) {
		return nil, errors.New("`samples` and `posteriorValues` must have the same length")
	}

	// Run preprocessing.
	// Sort along each axis.
	data := transpose(samples)
	sortingIndices := make([][]int, dimension)
	for i := 0; i < dimension; i++ {
		sortingIndices[i] = argSort(data[i])
	}

	// Copy posterior values.
	posterior := make([]float64, len(samples))
	copy(posterior, posteriorValues)

	return &voronoiBox{
		boundaries:     boundaries,
		sampleCount:    len(samples),
		data:           data,
		posterior:      posterior,
		sortingIndices: sortingIndices,
	}, nil
}

// TODO(p-nordmann): thread-safety.
type voronoiBox struct {
	boundaries     Boundaries
	sampleCount    int
	data           [][]float64
	posterior      []float64
	sortingIndices [][]int
}

func (box *voronoiBox) NewSampler(origin []float64) (Sampler, error) {
	dimension := box.boundaries.dimension()

	// Validate origin's dimension.
	if len(origin) != dimension {
		return nil, errors.New("`origin` has invalid dimension")
	}

	// Validate that origin is not out of bounds.
	for i, coordinate := range origin {
		if coordinate < box.boundaries.Infima[i] || box.boundaries.Suprema[i] < coordinate {
			return nil, errors.New("`origin` is out of bounds")
		}
	}

	// Build sampler: make buffers.
	sampler := &voronoiSampler{
		indexBuffer:            make([]int, box.sampleCount),
		intersectionBuffer:     make([]float64, box.sampleCount),
		box:                    box,
		squaredDistancesToAxis: make([]float64, box.sampleCount),
		position:               make([]float64, dimension),
	}
	copy(sampler.position, origin)

	// Compute samples' distances to last axis.
	for i := 0; i < dimension-1; i++ {
		coordinates := box.data[i]
		for k := range coordinates {
			sampler.squaredDistancesToAxis[k] += (origin[i] - coordinates[k]) * (origin[i] - coordinates[k])
		}
	}

	return sampler, nil
}

type voronoiSampler struct {
	// Buffers for computing minima.
	// We keep them alive so they are not garbage collected, for performance reasons.
	indexBuffer        []int
	intersectionBuffer []float64

	// Reference to the original data.
	box *voronoiBox

	// Last known position.
	// Is initialized with origin point and is updated during the walk.
	position []float64

	// Holds the squares of distances to the current axis to walk along.
	// It is passed as ordinates to the parabolas views for computing the minima.
	//
	// We keep them stored to avoid computing distances along each dimension each time we change direction.
	squaredDistancesToAxis []float64
}

// Sample walks to a new sample.
func (sampler *voronoiSampler) Sample() []float64 {
	dimension := sampler.box.boundaries.dimension()
	sample := make([]float64, dimension)
	copy(sample, sampler.position)
	for i := 0; i < dimension; i++ {
		// Make a step along axis d:

		// Update squared distances.
		previousI := i - 1
		if previousI < 0 {
			previousI += dimension
		}
		previousCoordinates := sampler.box.data[previousI]
		coordinates := sampler.box.data[i]
		for k := range sampler.squaredDistancesToAxis {
			sampler.squaredDistancesToAxis[k] += (sampler.position[previousI] - previousCoordinates[k]) * (sampler.position[previousI] - previousCoordinates[k])
			sampler.squaredDistancesToAxis[k] -= (sampler.position[i] - coordinates[k]) * (sampler.position[i] - coordinates[k])
		}

		// Compute minimal parabolas.
		parabolas(
			sampler.box.data[i], sampler.squaredDistancesToAxis, sampler.box.sortingIndices[i],
		).computeMinima(&(sampler.indexBuffer), &(sampler.intersectionBuffer))

		// Compute staircase density from minimal parabolas.
		infimum := sampler.box.boundaries.Infima[i]
		supremum := sampler.box.boundaries.Suprema[i]
		steps := []float64{infimum}
		stairs := []int{}
		for k := range sampler.indexBuffer {
			var nextIntersection float64
			if k >= len(sampler.intersectionBuffer) {
				nextIntersection = supremum
			} else {
				nextIntersection = sampler.intersectionBuffer[k]
			}
			// If we are out of bounds from the left, skip.
			if nextIntersection < infimum {
				continue
			}
			// Add a stair.
			stairs = append(stairs, sampler.indexBuffer[k])
			// If we are out of bounds from the right, terminate.
			if supremum <= nextIntersection {
				break
			}
			// Add the right bound of the stair.
			steps = append(steps, nextIntersection)
		}
		steps = append(steps, supremum)

		// Retrieve density values.
		density := make([]float64, len(stairs))
		for k, idx := range stairs {
			density[k] = sampler.box.posterior[sampler.box.sortingIndices[i][idx]]
		}

		// Sample coordinate.
		coordinate, _ := Staircase_naive(density, steps, false)

		// update sample and position
		sample[i] = coordinate
		sampler.position[i] = coordinate
	}
	return sample
}
