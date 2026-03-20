package math

import (
	"fmt"
	"math"
)

// CosineSimilarity calculates the cosine similarity between two vectors
// Returns a value between -1 and 1, where:
//
//	 1 = identical direction (same meaning)
//	 0 = orthogonal (unrelated)
//	-1 = opposite direction (opposite meaning)
func CosineSimilarity(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors must have same length: got %d and %d", len(a), len(b))
	}

	if len(a) == 0 {
		return 0, fmt.Errorf("cannot calculate similarity of empty vectors")
	}

	var dotProduct float64
	var sumA, sumB float64

	for i := range a {
		va := float64(a[i])
		vb := float64(b[i])

		dotProduct += va * vb
		sumA += va * va
		sumB += vb * vb
	}

	magnitudeA := math.Sqrt(sumA)
	magnitudeB := math.Sqrt(sumB)

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0, fmt.Errorf("cannot calculate cosine similarity with zero-magnitude vector")
	}

	similarity := dotProduct / (magnitudeA * magnitudeB)

	return float32(similarity), nil
}
