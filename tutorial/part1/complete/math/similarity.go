package math

import (
	"fmt"
	"math"
)

// CosineSimilarity calculates the cosine similarity between two vectors
// Returns a value between -1 and 1, where:
//   1 = identical direction (same meaning)
//   0 = orthogonal (unrelated)
//  -1 = opposite direction (opposite meaning)
func CosineSimilarity(a, b []float32) (float32, error) {
	// Validate inputs
	if len(a) != len(b) {
		return 0, fmt.Errorf("vectors must have same length: got %d and %d", len(a), len(b))
	}

	if len(a) == 0 {
		return 0, fmt.Errorf("cannot calculate similarity of empty vectors")
	}

	// Calculate dot product and magnitudes
	// Use float64 for intermediate calculations to avoid precision loss
	var dotProduct float64
	var sumA, sumB float64

	for i := range a {
		va := float64(a[i])
		vb := float64(b[i])

		dotProduct += va * vb
		sumA += va * va
		sumB += vb * vb
	}

	// Calculate magnitudes (Euclidean norms)
	magnitudeA := math.Sqrt(sumA)
	magnitudeB := math.Sqrt(sumB)

	// Check for zero magnitudes (avoid division by zero)
	if magnitudeA == 0 || magnitudeB == 0 {
		return 0, fmt.Errorf("cannot calculate cosine similarity with zero-magnitude vector")
	}

	// Calculate cosine similarity
	similarity := dotProduct / (magnitudeA * magnitudeB)

	// Convert back to float32 for consistency
	return float32(similarity), nil
}

// CalculateCentroid computes the centroid (average) of multiple vectors
// The centroid represents the "center" or "typical representative" of the group
func CalculateCentroid(vectors [][]float32) ([]float32, error) {
	// Validate input
	if len(vectors) == 0 {
		return nil, fmt.Errorf("cannot calculate centroid of empty vector set")
	}

	dimensions := len(vectors[0])

	// Check all vectors have same length
	for i, vec := range vectors {
		if len(vec) != dimensions {
			return nil, fmt.Errorf("vector %d has different length: expected %d, got %d", i, dimensions, len(vec))
		}
	}

	// Initialize centroid with zeros
	centroid := make([]float32, dimensions)
	numVectors := float32(len(vectors))

	// Sum each dimension across all vectors
	for i := 0; i < dimensions; i++ {
		var sum float32
		for _, vec := range vectors {
			sum += vec[i]
		}
		// Calculate average for this dimension
		centroid[i] = sum / numVectors
	}

	return centroid, nil
}
