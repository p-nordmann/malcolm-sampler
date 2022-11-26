package malcolms_test

import (
	"fmt"
	"testing"

	m "github.com/p-nordmann/malcolm-sampler"
)

// TODO(p-nordmann): small refactor to handle posterior normalization nicely.

// Tests that the API behaves correctly with edgy inputs.
//
// In particular we test here that it returns error on creating  sampler factories.
func TestEdgeCases(t *testing.T) {

	t.Run("should return error when malformed boundaries are provided", func(t *testing.T) {
		malformedBoundaries := m.Boundaries{
			Suprema: []float64{1, 2, 3},
			Infima:  []float64{0, 0},
		}
		samples := [][]float64{}
		posterior := []float64{}
		_, err := m.FromSamples(malformedBoundaries, samples, posterior)
		if err == nil {
			t.Error("got <nil> but expected error; boundaries have invalid dimensions")
		}
	})

	t.Run("should return error when invalid boundaries are provided", func(t *testing.T) {
		invalidBoundaries := m.Boundaries{
			Suprema: []float64{0, 0, 0},
			Infima:  []float64{1, 1, 1},
		}
		samples := [][]float64{{0.5, 0.5, 0.5}}
		posterior := []float64{1}
		_, err := m.FromSamples(invalidBoundaries, samples, posterior)
		if err == nil {
			t.Error("got <nil> but expected error; suprema are lower than infima")
		}
	})

	t.Run("should return error when dimension is 0", func(t *testing.T) {
		boundaries := m.Boundaries{
			Suprema: []float64{},
			Infima:  []float64{},
		}
		samples := [][]float64{{}}
		posterior := []float64{0}
		_, err := m.FromSamples(boundaries, samples, posterior)
		if err == nil {
			t.Error("got <nil> but expected error; dimension is 0")
		}
	})

	boundaries := m.Boundaries{
		Suprema: []float64{1, 1, 1},
		Infima:  []float64{0, 0, 0},
	}

	t.Run("should return error when samples don't have the correct dimension", func(t *testing.T) {
		samples := [][]float64{{0.5, 0.5}, {0.5, 0.5, 0.5}}
		posterior := []float64{1, 0}
		_, err := m.FromSamples(boundaries, samples, posterior)
		if err == nil {
			t.Error("got <nil> but expected error; samples don't have the correct dimension")
		}
	})

	t.Run("should return error when samples are out of bounds", func(t *testing.T) {
		samples := [][]float64{{0.5, 0.5, 0.5}, {0, 1.5, 0.5}}
		posterior := []float64{1, 0}
		_, err := m.FromSamples(boundaries, samples, posterior)
		if err == nil {
			t.Error("got <nil> but expected error; samples are out of bounds")
		}
	})

	t.Run("should return error when posterior and samples don't have the same length", func(t *testing.T) {
		samples := [][]float64{{0.5, 0.5, 0.5}}
		posterior := []float64{1, 2}
		_, err := m.FromSamples(boundaries, samples, posterior)
		if err == nil {
			t.Error("got <nil> but expected error; samples and posterior don't have the same length")
		}
	})

	t.Run("should return error when there is no sample", func(t *testing.T) {
		samples := [][]float64{}
		posterior := []float64{}
		_, err := m.FromSamples(boundaries, samples, posterior)
		if err == nil {
			t.Error("got <nil> but expected error; there is no sample")
		}
	})

	t.Run("NewSampler should return an error when `origin` is of wrong dimension", func(t *testing.T) {
		samples := [][]float64{{0.5, 0.5, 0.5}, {0, 0, 0.5}}
		posterior := []float64{1, 0}
		truePosterior, err := m.FromSamples(boundaries, samples, posterior)
		if err != nil {
			t.Errorf("returned an error: %v", err)
			return
		} else if truePosterior == nil {
			t.Error("returned <nil> but expected valid TruePosterior")
			return
		}
		_, err = truePosterior.NewSampler([]float64{0.5, 0.5})
		if err == nil {
			t.Error("got <nil> but expected error; origin is not of valid dimension")
		}
	})
}

// Helper for generating multiple samples.
//
// Checks that returned samples are of correct dimension and not out of bounds.
func walk(boundaries m.Boundaries, sampler m.Sampler, steps int) ([][]float64, error) {
	dimension := len(boundaries.Infima)
	generated := [][]float64{}
	for k := 0; k < steps; k++ {
		// Sample one point.
		sample := sampler.Sample()
		if len(sample) != dimension {
			return nil, fmt.Errorf(
				"invalid sampled dimension: got %d but expected %d", len(sample), dimension,
			)
		}
		for i := 0; i < dimension; i++ {
			if sample[i] < boundaries.Infima[i] || sample[i] > boundaries.Suprema[i] {
				return nil, fmt.Errorf(
					"sample out of bounds: %v out of bound on dimension %d", sample, i,
				)
			}
		}
		generated = append(generated, sample)
	}
	return generated, nil
}

// mockSampler loops over samples when asked to sample.
type mockSampler struct {
	samples   [][]float64
	callCount int
}

func (sampler *mockSampler) Sample() []float64 {
	defer func() { sampler.callCount++ }()
	return sampler.samples[sampler.callCount%len(sampler.samples)]
}

// Tests that walk function performs as expected.
func TestWalk(t *testing.T) {
	boundaries := m.Boundaries{
		Infima:  []float64{0, 0},
		Suprema: []float64{1, 1},
	}

	t.Run("should return an error when a sample's dimension doesn't match the boundaries", func(t *testing.T) {
		sampler := &mockSampler{
			samples: [][]float64{
				{0.5, 0.5},
				{0.25, 0.25},
				{0.5, 0.5, 0.5}, // Invalid dimension.
				{0.75, 0.75},
			},
		}
		_, err := walk(boundaries, sampler, 10)
		if err == nil {
			t.Error("returned <nil> but expected error")
		}
	})

	t.Run("should return an error when a sample is out of bounds", func(t *testing.T) {
		sampler := &mockSampler{
			samples: [][]float64{
				{0.5, 0.5},
				{0.25, 0.25},
				{0.5, 1.5}, // Out of bounds.
				{0.75, 0.75},
			},
		}
		_, err := walk(boundaries, sampler, 10)
		if err == nil {
			t.Error("returned <nil> but expected error")
		}
	})

	t.Run("should return the correct number of points", func(t *testing.T) {
		sampler := &mockSampler{
			samples: [][]float64{
				{0.5, 0.5},
				{0.25, 0.25},
			},
		}
		got, err := walk(boundaries, sampler, 5)
		if err != nil {
			t.Errorf("returned error: %v", err)
		}
		want := [][]float64{
			{0.5, 0.5},
			{0.25, 0.25},
			{0.5, 0.5},
			{0.25, 0.25},
			{0.5, 0.5},
		}
		for k := range want {
			if got[k][0] != want[k][0] || got[k][1] != want[k][1] {
				t.Errorf("%dth sample: got %v but expected %v", k, got[k], want[k])
			}
		}
	})
}

// chiSquared computes the chi-squared value of the observations against the expectations.
//
// Example for a dice:
//
//	We perform 100 throws and get 24 ones, 16 twos, 12 threes, 14 fours, 19 fives, 15 sixes.
//	observations := []int{24, 16, 12, 14, 19, 15}
//	expectations := []float64{1.0 / 6.0, 1.0 / 6.0, 1.0 / 6.0, 1.0 / 6.0, 1.0 / 6.0, 1.0 / 6.0}
//	chiSquared(observations, expectations)
//	> 5.48
func chiSquared(observations []int, expectations []float64) float64 {
	var total int
	for _, count := range observations {
		total += count
	}
	var chi2 float64
	for k, expected := range expectations {
		expectedCount := float64(total) * expected
		diff := (float64(observations[k]) - expectedCount)
		chi2 += diff * diff / expectedCount
	}
	return chi2
}

// Tests that we can sample approximations from usual 1-D distributions.
func TestSamplingUsualDistributions1D(t *testing.T) {
	// TODO(p-nordmann): implementation.
}

// Creates samples at the center of the smaller cubes obtains in cutting each side of the cube in half.
//
// Dimension 2 example:
//
//	0   1   2
//	+---+---+ 2
//	| x | x |
//	+---+---+ 1
//	| x | x |
//	+---+---+ 0
//	leading to samples {0.5, 0.5}, {0.5, 1.5}, {1.5, 0.5}, {1.5, 1.5}.
func createCubeSamples(dimension int) [][]float64 {
	// We want n = 2^d samples:
	n := 1 << dimension
	samples := make([][]float64, n)
	// We will use the writing of integers in base 2 to determine whether to use 0.5 or 1.5 coordinates.
	for k := 0; k < n; k++ {
		samples[k] = make([]float64, dimension)
		for i := 0; i < dimension; i++ {
			// Mask for i-th bit is 1 << i.
			if k&(1<<i) == 0 {
				samples[k][i] = 0.5
			} else {
				samples[k][i] = 1.5
			}
		}
	}
	return samples
}

// Tests that createCubeSamples function performs as expected.
func TestCreateCubeSamples(t *testing.T) {

	// compareSamples returns whether the sets of samples are the same.
	compareSamples := func(got, want [][]float64) bool {
		if len(got) != len(want) {
			return false
		}
		for k := range want {
			if len(got[k]) != len(want[k]) {
				return false
			}
			for i := range want[k] {
				if got[k][i] != want[k][i] {
					return false
				}
			}
		}
		return true
	}

	t.Run("dimension 1", func(t *testing.T) {
		want := [][]float64{
			{0.5}, {1.5},
		}
		got := createCubeSamples(1)
		if !compareSamples(got, want) {
			t.Errorf("invalid cube samples for dimension 1: %v", got)
		}
	})

	t.Run("dimension 2", func(t *testing.T) {
		want := [][]float64{
			{0.5, 0.5}, {1.5, 0.5}, {0.5, 1.5}, {1.5, 1.5},
		}
		got := createCubeSamples(2)
		if !compareSamples(got, want) {
			t.Errorf("invalid cube samples for dimension 2: %v", got)
		}
	})

	t.Run("dimension 3", func(t *testing.T) {
		want := [][]float64{
			{0.5, 0.5, 0.5}, {1.5, 0.5, 0.5}, {0.5, 1.5, 0.5}, {1.5, 1.5, 0.5},
			{0.5, 0.5, 1.5}, {1.5, 0.5, 1.5}, {0.5, 1.5, 1.5}, {1.5, 1.5, 1.5},
		}
		got := createCubeSamples(3)
		if !compareSamples(got, want) {
			t.Errorf("invalid cube samples for dimension 3: %v", got)
		}
	})
}

// validateDistribution runs chi-squared test on the generated samples against the expected
// distribution.
//
// The expected distribution is uniform across each unique cube. We obtain a discrete distribution
// by gathering samples that fall in the same small cube. Each cube is expected to contain a number
// of samples proportional to `posterior[k]` (volume is 1) where k is the index of the cube as
// generated by `createCubeSamples`.
//
// NB: this does not test that distribution is uniform inside one cube. For this, a custom test
// should generate a set of samples from a lower number of posterior samples.
//
// NB: the ordering of samples generated by `createCubeSamples` matters!
//
// TODO(p-nordmann): build the test the other way around; so that the fundamental hypothesis is
// rejecting the distribution.
func validateCubeDistribution(dimension int, samples [][]float64, posterior []float64) bool {

	// Count samples in each small cube.
	counts := make([]int, 1<<dimension)
	for _, sample := range samples {
		var cubeIndex int
		var mask int = 1
		for i := 0; i < dimension; i++ {
			if sample[i] >= 1 {
				cubeIndex |= mask
			}
			mask <<= 1
		}
		counts[cubeIndex]++
	}

	// Normalize expected posterior.
	var posteriorSum float64
	for _, value := range posterior {
		posteriorSum += value
	}
	normalizedPosterior := make([]float64, len(posterior))
	for k := range posterior {
		normalizedPosterior[k] = posterior[k] / posteriorSum
	}

	// Compare to expected distribution using chi-squared and table.
	chi2 := chiSquared(counts, normalizedPosterior)
	table5 := []float64{
		3.84146, 5.99146, 7.81473, 9.48773, 11.07050, 12.59159, 14.06714, 15.50731, 16.91898,
		18.30704, 19.67514, 21.02607, 22.36203, 23.68479, 24.99579, 26.29623,
	}
	return chi2 < table5[(1<<dimension)-2]
}

type samplingTestCase struct {
	description string
	boundaries  m.Boundaries
	trueSamples struct {
		samples   [][]float64
		posterior []float64
	}
	origin            []float64
	expectedPosterior []float64
}

// Tests that the API supports different dimensions for the parameter space.
func TestSamplingWithVariableDimensions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test because time consuming")
	}

	const maxDimension = 4 // Limited by chi-squared quantile table; need to go up to 2^d-1.
	const sampleCountRequest = 10000
	var testCases []samplingTestCase

	for dimension := 1; dimension <= maxDimension; dimension++ {
		// Build boundaries.
		boundaries := m.Boundaries{
			Suprema: make([]float64, dimension),
			Infima:  make([]float64, dimension),
		}
		for k := 0; k < dimension; k++ {
			boundaries.Suprema[k] = 2
			boundaries.Infima[k] = 0
		}
		// Build 2^d samples.
		samples := createCubeSamples(dimension)
		expectedPosterior := make([]float64, len(samples))
		// Single central point.
		point := make([]float64, dimension)
		for i := 0; i < dimension; i++ {
			point[i] = 1
		}

		// With 2^d true samples.
		for k := range expectedPosterior {
			expectedPosterior[k] = float64(k + 1)
		}
		testCases = append(testCases, samplingTestCase{
			description: "should sample according to staircase density",
			boundaries:  boundaries,
			trueSamples: struct {
				samples   [][]float64
				posterior []float64
			}{
				samples:   samples,
				posterior: expectedPosterior,
			},
			origin:            point,
			expectedPosterior: expectedPosterior,
		})

		// With 1 true sample.
		expectedPosterior = make([]float64, len(samples))
		for k := range expectedPosterior {
			expectedPosterior[k] = 1 // uniform accross the cube.
		}
		testCases = append(testCases, samplingTestCase{
			description: "should sample uniformly",
			boundaries:  boundaries, trueSamples: struct {
				samples   [][]float64
				posterior []float64
			}{
				samples:   [][]float64{point},
				posterior: []float64{1},
			},
			origin:            point,
			expectedPosterior: expectedPosterior,
		})
	}
	for k, testCase := range testCases {
		t.Run(fmt.Sprintf("case #%d: %s", k, testCase.description), func(t *testing.T) {

			// Create factory.
			truePosterior, err := m.FromSamples(
				testCase.boundaries, testCase.trueSamples.samples, testCase.trueSamples.posterior,
			)
			if err != nil {
				t.Errorf("returned error: %v", err)
				return
			} else if truePosterior == nil {
				t.Error("returned <nil> but expected valid TruePosterior")
				return
			}

			// Create sampler.
			sampler, err := truePosterior.NewSampler(testCase.origin)
			if err != nil {
				t.Errorf("returned error: %v", err)
				return
			} else if sampler == nil {
				t.Error("returned <nil> but expected valid Sampler")
				return
			}

			// Walk using sampler and validate results.
			dimension := len(testCase.boundaries.Infima)
			sampled, err := walk(
				testCase.boundaries, sampler, sampleCountRequest,
			)
			if err != nil {
				t.Errorf("error during walk: %v", err)
				return
			} else if !validateCubeDistribution(
				dimension, sampled, testCase.expectedPosterior,
			) {
				t.Error("sampled distribution does not match expected distribution")
				return
			}
		})
	}
}
