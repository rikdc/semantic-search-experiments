package math

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// centerMatrix subtracts the mean of each dimension from all vectors.
// PCA requires zero-mean data; skipping this step produces axes that describe
// distance from the origin rather than directions of variance.
func centerMatrix(data [][]float32) [][]float64 {
	n := len(data)
	d := len(data[0])
	means := make([]float64, d)

	for _, v := range data {
		for j, val := range v {
			means[j] += float64(val)
		}
	}
	for j := range means {
		means[j] /= float64(n)
	}

	centered := make([][]float64, n)
	for i, v := range data {
		centered[i] = make([]float64, d)
		for j, val := range v {
			centered[i][j] = float64(val) - means[j]
		}
	}
	return centered
}

// ProjectTo2D reduces a set of high-dimensional vectors to 2D using PCA.
// All vectors in data must have the same length.
//
// If you want to include a query vector in the plot, append it to data
// before calling this function. The query must share the same PCA axes as
// the training points to be plotted in the same coordinate space.
func ProjectTo2D(data [][]float32) ([][2]float64, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("cannot project empty data")
	}
	if len(data) < 2 {
		return nil, fmt.Errorf("need at least 2 vectors for PCA")
	}

	centered := centerMatrix(data)
	n := len(centered)
	d := len(centered[0])

	flat := make([]float64, n*d)
	for i, row := range centered {
		copy(flat[i*d:], row)
	}

	m := mat.NewDense(n, d, flat)
	var svd mat.SVD
	if ok := svd.Factorize(m, mat.SVDThin); !ok {
		return nil, fmt.Errorf("SVD factorization failed")
	}

	// svd.VTo returns V (d×min(n,d)), not V-transposed.
	// The k-th principal component is the k-th column of V, so we access V.At(j, k).
	var v mat.Dense
	svd.VTo(&v)

	_, nComponents := v.Dims()
	if nComponents < 2 {
		return nil, fmt.Errorf("not enough singular values for 2D projection (got %d)", nComponents)
	}

	points := make([][2]float64, n)
	for i := range points {
		for k := 0; k < 2; k++ {
			for j := 0; j < d; j++ {
				points[i][k] += centered[i][j] * v.At(j, k)
			}
		}
	}
	return points, nil
}
