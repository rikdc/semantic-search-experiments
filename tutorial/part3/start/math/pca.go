package math

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// centerMatrix subtracts the mean of each dimension from all vectors.
// PCA requires zero-mean data; skipping this step produces axes that describe
// distance from the origin rather than directions of variance.
//
// TODO: implement centerMatrix
// 1. Compute the mean of each dimension across all n vectors.
// 2. Subtract each dimension's mean from every vector.
// Hint: accumulate sums in a float64 slice, divide by n, then subtract.
func centerMatrix(data [][]float32) [][]float64 {
	panic("not implemented")
}

// ProjectTo2D reduces a set of high-dimensional vectors to 2D using PCA.
// All vectors in data must have the same length.
//
// If you want to include a query vector in the plot, append it to data
// before calling this function. The query must share the same PCA axes as
// the training points to be plotted in the same coordinate space.
//
// TODO: implement ProjectTo2D
// 1. Call centerMatrix(data) to get zero-mean float64 data.
// 2. Flatten the centered data into a single []float64 slice (row-major).
// 3. Build a *mat.Dense from the flat slice: mat.NewDense(n, d, flat).
// 4. Factorize with svd.Factorize(m, mat.SVDThin) — check the bool return.
// 5. Call svd.VTo(&v) to get V, a d×min(n,d) matrix.
// 6. Project: for each point i and component k in 0..1,
//    points[i][k] = sum over j of centered[i][j] * v.At(j, k)
// Hint: v.At(j, k) gives the j-th element of the k-th principal component (column k of V).
//       Do NOT use v.At(k, j) — that transposes the indices and panics on tall matrices.
func ProjectTo2D(data [][]float32) ([][2]float64, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("cannot project empty data")
	}
	if len(data) < 2 {
		return nil, fmt.Errorf("need at least 2 vectors for PCA")
	}

	// Suppress unused import — remove this line once you use mat below.
	_ = mat.SVDThin

	panic("not implemented")
}
