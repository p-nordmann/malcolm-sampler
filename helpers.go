package malcolms

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// argSort returns a list of indices in sorted order.
//
// Uses stable sorting.
func argSort(floats []float64) []int {
	if len(floats) == 0 {
		return nil
	}
	indices := make([]int, len(floats))
	for k := range indices {
		indices[k] = k
	}
	sort.SliceStable(indices, func(i, j int) bool {
		return floats[indices[i]] < floats[indices[j]]
	})
	return indices
}

// transpose returns the transposition of `matrix`.
//
// It is useful to get from row-major to column-major indexing.
//
// Expects all lines to have the same dimension and height to be >= 1.
func transpose(matrix [][]float64) [][]float64 {
	height, width := len(matrix), len(matrix[0])
	transposedMatrix := make([][]float64, width)
	for k := range transposedMatrix {
		transposedMatrix[k] = make([]float64, height)
		for l := range transposedMatrix[k] {
			transposedMatrix[k][l] = matrix[l][k]
		}
	}
	return transposedMatrix
}

// Generates a random value according to a 1D staircase density of probability.
// Uses the naive rejection method.
// TODO: safeguards for infinite loop
//
// TODO(p-nordmann): piece of code from old codebase. To be refactored.
func Staircase_naive(density, steps []float64, islog bool) (float64, int) {
	// retrieve n
	n := len(density)
	if n == 0 {
		panic(ErrorVoidStaircase(0))
	} else if len(steps) != n+1 {
		panic(ErrorStaircaseMismatch{n, len(steps)})
	}

	// retrieve maximum of density
	dmax := density[0]
	for _, d := range density {
		if d > dmax {
			dmax = d
		}
	}

	// rejection loop
	count := 0
	for {
		count++
		x := rand.Float64()*(steps[n]-steps[0]) + steps[0]
		// use bisection to locate x
		a, b := 0, n-1
		t := (a + b) / 2
		for a < b {
			if steps[t+1] <= x {
				a = t + 1
			} else {
				b = t
			}
			t = (a + b) / 2
		}
		// rejection method
		r := rand.Float64()
		if islog {
			r = math.Log(r) + dmax
		} else {
			r *= dmax
		}
		if r <= density[t] {
			return x, t
		}
	}
}

type ErrorVoidStaircase int

func (e ErrorVoidStaircase) Error() string {
	return "Void staircase given for generation."
}

type ErrorStaircaseMismatch [2]int

func (e ErrorStaircaseMismatch) Error() string {
	return fmt.Sprintf("Density and steps length mismatch for given staircase (%d,%d).", e[0], e[1])
}
