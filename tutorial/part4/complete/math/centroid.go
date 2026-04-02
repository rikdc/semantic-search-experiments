package math

import "fmt"

// CalculateCentroid computes the centroid (element-wise mean) of a set of vectors.
// The centroid represents the "average location" in embedding space for the group —
// used here as the reference point for Normal transactions.
func CalculateCentroid(vectors [][]float32) ([]float32, error) {
	if len(vectors) == 0 {
		return nil, fmt.Errorf("cannot calculate centroid of empty vector set")
	}

	dims := len(vectors[0])
	for i, vec := range vectors {
		if len(vec) != dims {
			return nil, fmt.Errorf("vector %d has different length: expected %d, got %d", i, dims, len(vec))
		}
	}

	centroid := make([]float32, dims)
	n := float32(len(vectors))

	for i := 0; i < dims; i++ {
		var sum float32
		for _, vec := range vectors {
			sum += vec[i]
		}
		centroid[i] = sum / n
	}

	return centroid, nil
}
