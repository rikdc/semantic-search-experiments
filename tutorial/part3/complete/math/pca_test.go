package math

import (
	"math"
	"testing"
)

func TestCenterMatrix_ZeroMean(t *testing.T) {
	data := [][]float32{
		{1, 2, 3},
		{3, 4, 5},
		{5, 6, 7},
	}
	centered := centerMatrix(data)

	// Each dimension should sum to ~0 after centering
	d := len(data[0])
	for j := 0; j < d; j++ {
		var sum float64
		for i := range centered {
			sum += centered[i][j]
		}
		if math.Abs(sum) > 1e-10 {
			t.Errorf("dimension %d not zero-mean after centering: sum = %v", j, sum)
		}
	}
}

func TestCenterMatrix_KnownValues(t *testing.T) {
	data := [][]float32{
		{1, 4},
		{3, 6},
	}
	// means: [2, 5] → centered: [[-1, -1], [1, 1]]
	centered := centerMatrix(data)

	expected := [][]float64{{-1, -1}, {1, 1}}
	for i, row := range centered {
		for j, val := range row {
			if math.Abs(val-expected[i][j]) > 1e-10 {
				t.Errorf("centered[%d][%d] = %v, want %v", i, j, val, expected[i][j])
			}
		}
	}
}

func TestProjectTo2D_OutputDimensions(t *testing.T) {
	data := [][]float32{
		{1, 0, 0},
		{2, 0, 0},
		{0, 1, 0},
		{0, 2, 0},
		{0, 0, 1},
	}
	points, err := ProjectTo2D(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(points) != len(data) {
		t.Errorf("got %d points, want %d", len(points), len(data))
	}
	for i, p := range points {
		if math.IsNaN(p[0]) || math.IsNaN(p[1]) {
			t.Errorf("point %d has NaN coordinates", i)
		}
	}
}

func TestProjectTo2D_PreservesClusterSeparation(t *testing.T) {
	// Two clearly separated clusters in 3D.
	// After projection to 2D, the between-cluster distance should remain large.
	clusterA := [][]float32{
		{10, 0, 0},
		{10.1, 0.1, -0.1},
		{9.9, -0.1, 0.1},
	}
	clusterB := [][]float32{
		{0, 10, 0},
		{0.1, 10.1, -0.1},
		{-0.1, 9.9, 0.1},
	}

	all := append(clusterA, clusterB...)
	points, err := ProjectTo2D(all)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Compute 2D centroids of each cluster
	centA := [2]float64{}
	for i := 0; i < 3; i++ {
		centA[0] += points[i][0]
		centA[1] += points[i][1]
	}
	centA[0] /= 3
	centA[1] /= 3

	centB := [2]float64{}
	for i := 3; i < 6; i++ {
		centB[0] += points[i][0]
		centB[1] += points[i][1]
	}
	centB[0] /= 3
	centB[1] /= 3

	betweenDist := math.Sqrt(math.Pow(centA[0]-centB[0], 2) + math.Pow(centA[1]-centB[1], 2))
	if betweenDist < 5.0 {
		t.Errorf("clusters should be well-separated in 2D, centroid distance = %v", betweenDist)
	}
}

func TestProjectTo2D_QueryAppendedLast(t *testing.T) {
	// Appending the query before PCA should place it at the last index,
	// in the same coordinate space as the training points.
	training := [][]float32{
		{1, 0, 0},
		{2, 0, 0},
		{0, 1, 0},
		{0, 2, 0},
	}
	query := []float32{1.5, 0.5, 0}

	allVectors := append(training, query)
	points, err := ProjectTo2D(allVectors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(points) != len(allVectors) {
		t.Errorf("expected %d points, got %d", len(allVectors), len(points))
	}

	queryPoint := points[len(allVectors)-1]
	if math.IsNaN(queryPoint[0]) || math.IsNaN(queryPoint[1]) {
		t.Errorf("query point has NaN coordinates")
	}
}

func TestProjectTo2D_EmptyInput(t *testing.T) {
	_, err := ProjectTo2D([][]float32{})
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestProjectTo2D_SingleVector(t *testing.T) {
	_, err := ProjectTo2D([][]float32{{1, 2, 3}})
	if err == nil {
		t.Fatal("expected error for single vector, got nil")
	}
}
