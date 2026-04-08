package math

// CalculateMeanStdDev computes the mean and population standard deviation
// of a slice of similarity scores.
//
// Both values are returned as float64 for precision — the standard deviation
// of a tight Normal cluster can be small (≈0.01), and float32 arithmetic
// in the denominator of z-score calculations introduces meaningful error.
//
// Target behaviour:
//   - Empty slice → return (0, 0)
//   - All identical values → stddev = 0
//   - sims = {0.80, 0.84, 0.86, 0.90} → mean ≈ 0.85, stddev ≈ 0.036
func CalculateMeanStdDev(sims []float32) (mean, stddev float64) {
	// TODO: implement CalculateMeanStdDev
	// Step 1: sum all values into mean, then divide by len(sims)
	// Step 2: for each value, accumulate (float64(s) - mean)^2 into stddev
	// Step 3: stddev = math.Sqrt(stddev / float64(len(sims)))
	panic("not implemented")
}
