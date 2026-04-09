// Package math provides vector arithmetic functions shared across all parts
// of the It's Just Vectors tutorial series.
package math

import (
	"fmt"
	"math"
)

// CosineSimilarity returns the cosine similarity between two vectors a and b.
// The result ranges from -1.0 (opposite) to 1.0 (identical direction).
// Intermediate calculations use float64 to avoid rounding error in the
// accumulation loop; the diagonal of a similarity matrix should be exactly 1.0.
func CosineSimilarity(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vector length mismatch: %d vs %d", len(a), len(b))
	}
	if len(a) == 0 {
		return 0, fmt.Errorf("vectors must be non-empty")
	}

	var dot, magA, magB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		magA += float64(a[i]) * float64(a[i])
		magB += float64(b[i]) * float64(b[i])
	}

	magA = math.Sqrt(magA)
	magB = math.Sqrt(magB)

	if magA == 0 || magB == 0 {
		return 0, fmt.Errorf("zero-magnitude vector")
	}

	return float32(dot / (magA * magB)), nil
}

// CalculateMeanStdDev returns the mean and population standard deviation of a
// slice of similarity scores. Uses a two-pass approach for numerical stability:
// when similarities are tightly clustered, single-pass Welford can lose precision
// in the variance accumulation step.
func CalculateMeanStdDev(sims []float32) (mean, stddev float64) {
	if len(sims) == 0 {
		return 0, 0
	}

	n := float64(len(sims))
	for _, s := range sims {
		mean += float64(s)
	}
	mean /= n

	for _, s := range sims {
		diff := float64(s) - mean
		stddev += diff * diff
	}
	stddev = math.Sqrt(stddev / n)
	return mean, stddev
}
