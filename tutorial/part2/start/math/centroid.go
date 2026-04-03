package math

import "fmt"

// CalculateCentroid computes the centroid (average) of multiple vectors
// The centroid represents the "center" or "typical representative" of the group
func CalculateCentroid(vectors [][]float32) ([]float32, error) {
	if len(vectors) == 0 {
		return nil, fmt.Errorf("cannot calculate centroid of empty vector set")
	}

	dimensions := len(vectors[0])

	for i, vec := range vectors {
		if len(vec) != dimensions {
			return nil, fmt.Errorf("vector %d has different length: expected %d, got %d", i, dimensions, len(vec))
		}
	}

	centroid := make([]float32, dimensions)
	numVectors := float64(len(vectors))

	for i := 0; i < dimensions; i++ {
		var sum float64
		for _, vec := range vectors {
			sum += float64(vec[i])
		}
		centroid[i] = float32(sum / numVectors)
	}

	return centroid, nil
}
