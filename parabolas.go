package malcolms

// TODO(p-nordmann): separate package with implementation tests?
//	the logic here is far from trivial, it would make sense to have its own place and tests.

// parabola represents a 1D unit parabola.
type parabola struct {
	absciss, ordinate float64
}

// intersection returns the absciss of the intersection between two unit parabolas.
//
// Panics when comparing identical parabolas because intersection is infinite.
//
// TODO(p-nordmann): what about a branchless function? Find whether NaN is acceptable.
func intersection(p1, p2 parabola) float64 {
	if p1.absciss == p2.absciss {
		if p1.ordinate == p2.ordinate {
			panic("identical parabolas; infinite intersection")
		}
		panic("parallel parabolas; no intersection")
	}
	return 0.5 * (p1.ordinate - p2.ordinate + p1.absciss*p1.absciss - p2.absciss*p2.absciss) / (p1.absciss - p2.absciss)
}

// parabolasView encapsulates the logic to deal with a sequence of 1D unit parabolas.
type parabolasView struct {
	abscisses, ordinates []float64
	sortingIndex         []int
}

func parabolas(abscisses, ordinates []float64, sortingIndex []int) parabolasView {
	return parabolasView{
		abscisses:    abscisses,
		ordinates:    ordinates,
		sortingIndex: sortingIndex,
	}
}

// len returns the length of the sequence.
func (view parabolasView) len() int {
	return len(view.sortingIndex)
}

// get returns the k-th parabola of the sequence.
func (view parabolasView) get(k int) parabola {
	return parabola{
		absciss:  view.abscisses[view.sortingIndex[k]],
		ordinate: view.ordinates[view.sortingIndex[k]],
	}
}

// computeMinima finds the list of minimal parabolas within the sequence.
//
// Stores results in input buffers; buffers are expected to have big-enough capacity.
//
// TODO(p-nordmann): error management + handling parallel parabolas.
//
// TODO(p-nordmann): do not put negligible parabolas in results. Verify that this is guaranteed.
func (view parabolasView) computeMinima(indexBuffer *[]int, intersectionBuffer *[]float64) {
	*indexBuffer = (*indexBuffer)[:0]
	*intersectionBuffer = (*intersectionBuffer)[:0]

	for k := 0; k < view.len(); k++ {
		for len(*indexBuffer) > 0 {
			lastIndex := (*indexBuffer)[len(*indexBuffer)-1]

			// Edge case: previous absciss is the same as current absciss.
			// Only keep current parabola if strictly below previous parabola.
			if view.get(lastIndex).absciss == view.get(k).absciss {
				if view.get(lastIndex).ordinate <= view.get(k).ordinate {
					break
				}
				*indexBuffer = (*indexBuffer)[0 : len(*indexBuffer)-1]
				if len(*intersectionBuffer) > 0 {
					*intersectionBuffer = (*intersectionBuffer)[0 : len(*intersectionBuffer)-1]
				}
				continue
			}

			intersection := intersection(view.get(k), view.get(lastIndex))
			if len(*indexBuffer) == 1 {
				*intersectionBuffer = append(*intersectionBuffer, intersection)
				*indexBuffer = append(*indexBuffer, k)
				break
			}
			lastIntersection := (*intersectionBuffer)[len(*intersectionBuffer)-1]
			if lastIntersection < intersection {
				*intersectionBuffer = append(*intersectionBuffer, intersection)
				*indexBuffer = append(*indexBuffer, k)
				break
			}
			*indexBuffer = (*indexBuffer)[0 : len(*indexBuffer)-1]
			*intersectionBuffer = (*intersectionBuffer)[0 : len(*intersectionBuffer)-1]
		}

		if len(*indexBuffer) == 0 {
			*indexBuffer = append(*indexBuffer, k)
		}
	}
}
