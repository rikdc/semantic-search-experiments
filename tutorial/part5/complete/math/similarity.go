package math

import (
	"fmt"
	"math"
)

// CosineSimilarity calculates the cosine similarity between two vectors.
// Returns a value between -1 and 1, where:
//
//	 1 = identical direction (same meaning)
//	 0 = orthogonal (unrelated)
//	-1 = opposite direction (opposite meaning)
//
// Uses float64 for intermediate accumulation to avoid rounding error —
// a vector compared to itself must return exactly 1.0000.
func CosineSimilarity(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors must have same length: got %d and %d", len(a), len(b))
	}
	if len(a) == 0 {
		return 0, fmt.Errorf("cannot calculate similarity of empty vectors")
	}

	var dotProduct, sumA, sumB float64
	for i := range a {
		va := float64(a[i])
		vb := float64(b[i])
		dotProduct += va * vb
		sumA += va * va
		sumB += vb * vb
	}

	magA := math.Sqrt(sumA)
	magB := math.Sqrt(sumB)
	if magA == 0 || magB == 0 {
		return 0, fmt.Errorf("cannot calculate cosine similarity with zero-magnitude vector")
	}

	return float32(dotProduct / (magA * magB)), nil
}
