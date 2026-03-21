package math

import (
	"testing"
)

func TestCalculateCentroid_KnownVectors(t *testing.T) {
	vectors := [][]float32{
		{2, 4},
		{4, 6},
		{6, 8},
	}

	centroid, err := CalculateCentroid(vectors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Expected: [4, 6]
	if len(centroid) != 2 {
		t.Fatalf("expected 2 dimensions, got %d", len(centroid))
	}
	if centroid[0] != 4.0 {
		t.Errorf("centroid[0] = %f, want 4.0", centroid[0])
	}
	if centroid[1] != 6.0 {
		t.Errorf("centroid[1] = %f, want 6.0", centroid[1])
	}
}

func TestCalculateCentroid_EmptyInput(t *testing.T) {
	_, err := CalculateCentroid([][]float32{})
	if err == nil {
		t.Fatal("expected error for empty input, got nil")
	}
}

func TestCalculateCentroid_MismatchedDimensions(t *testing.T) {
	vectors := [][]float32{
		{1, 2, 3},
		{4, 5},
	}

	_, err := CalculateCentroid(vectors)
	if err == nil {
		t.Fatal("expected error for mismatched dimensions, got nil")
	}
}

func TestNearestCategory_ClearWinner(t *testing.T) {
	// Simulate category centroids with known vectors
	coffeeCentroid := []float32{0.9, 0.1, 0.0}
	groceryCentroid := []float32{0.1, 0.9, 0.0}
	gasCentroid := []float32{0.0, 0.1, 0.9}

	// Query that clearly resembles coffee
	query := []float32{0.85, 0.15, 0.05}

	centroids := map[string][]float32{
		"coffee":  coffeeCentroid,
		"grocery": groceryCentroid,
		"gas":     gasCentroid,
	}

	bestCategory := ""
	var bestSim float32

	for category, centroid := range centroids {
		sim, err := CosineSimilarity(query, centroid)
		if err != nil {
			t.Fatalf("unexpected error for category %s: %v", category, err)
		}
		if sim > bestSim {
			bestSim = sim
			bestCategory = category
		}
	}

	if bestCategory != "coffee" {
		t.Errorf("expected category 'coffee', got '%s' (similarity: %.4f)", bestCategory, bestSim)
	}
}

func TestNearestCategory_Tie(t *testing.T) {
	// Two centroids equidistant from the query
	centroidA := []float32{1, 0}
	centroidB := []float32{0, 1}

	// Query at 45 degrees — equidistant from both
	query := []float32{1, 1}

	simA, err := CosineSimilarity(query, centroidA)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	simB, err := CosineSimilarity(query, centroidB)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both similarities should be equal (cosine of 45 degrees)
	diff := simA - simB
	if diff < 0 {
		diff = -diff
	}

	if diff > 1e-6 {
		t.Errorf("expected equal similarities for tie, got %.6f and %.6f (diff: %.6f)", simA, simB, diff)
	}
}
