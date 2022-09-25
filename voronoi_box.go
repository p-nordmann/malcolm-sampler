package malcolms

// FromSamples runs preprocessing on samples and posterior values and returns a structure that
// implements TruePosterior interface.
//
// Dimensions are expected to be consistent between all parameters.
func FromSamples(
	boundaries Boundaries, samples [][]float64, posteriorValues []float64,
) (TruePosterior, error) {
	return nil, nil
}
