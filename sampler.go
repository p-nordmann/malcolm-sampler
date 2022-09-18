package malcolms

// Boundaries encapsulate the bounding box of the inversion problem.
type Boundaries struct {
	// Suprema and Infima holds the upper and lower limit values of the inversion problem.
	// They must have the same dimension.
	Suprema, Infima []float64
}

// TruePosterior encapsulate the samples that have a true posterior value.
//
// It provides all the necessary preprocessing and acts as a sampler factory.
type TruePosterior interface {
	NewSampler() Sampler
}

// Sampler is a thread-safe generator of samples.
//
// It refers to the TruePosterior that created it for the true posterior values and necessary
// information.
type Sampler interface {
	Sample() []float64
}
