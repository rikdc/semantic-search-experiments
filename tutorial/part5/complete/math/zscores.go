package math

import "math"

// CalculateMeanStdDev computes the mean and population standard deviation
// of a slice of similarity scores.
//
// Both values are returned as float64 for precision — the standard deviation
// of a tight Normal cluster can be small (≈0.01), and float32 arithmetic
// in the denominator of z-score calculations introduces meaningful error.
//
// Uses population standard deviation (divide by n), not sample (divide by n-1).
// For small training sets the difference is minor; for larger datasets prefer
// the sample form.
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
