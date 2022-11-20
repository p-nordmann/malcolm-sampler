package malcolms

// Boundaries encapsulate the bounding box of the inversion problem.
type Boundaries struct {
	// Suprema and Infima holds the upper and lower limit values of the inversion problem.
	// They must have the same dimension.
	Suprema, Infima []float64
}

func (boundaries Boundaries) dimension() int {
	return len(boundaries.Infima)
}

// SamplerFactory encapsulate the samples that have a true posterior value.
//
// It provides all the necessary preprocessing and acts as a sampler factory.
type SamplerFactory interface {
	// NewSampler returns a Sampler ready to walk from `origin`.
	//
	// It expects `origin` to be of correct dimension.
	NewSampler(origin []float64) (Sampler, error)
}

// Sampler is a thread-safe generator of samples.
//
// It refers to its parent factory to get the data it needs for sampling.
type Sampler interface {
	// Sample walks and returns a new sample.
	Sample() []float64
}
