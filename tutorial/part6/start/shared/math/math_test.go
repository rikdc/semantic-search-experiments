package math_test

import (
	"math"
	"testing"

	vmath "part6/shared/math"
)

func TestCosineSimilarity_IdenticalVectors(t *testing.T) {
	a := []float32{1.0, 2.0, 3.0}
	sim, err := vmath.CosineSimilarity(a, a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(float64(sim)-1.0) > 1e-6 {
		t.Errorf("identical vectors: expected similarity 1.0, got %.6f", sim)
	}
}

func TestCosineSimilarity_OppositeVectors(t *testing.T) {
	a := []float32{1.0, 0.0, 0.0}
	b := []float32{-1.0, 0.0, 0.0}
	sim, err := vmath.CosineSimilarity(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(float64(sim)+1.0) > 1e-6 {
		t.Errorf("opposite vectors: expected similarity -1.0, got %.6f", sim)
	}
}

func TestCosineSimilarity_ScaledVectors(t *testing.T) {
	// B is twice A — same direction, so similarity must be 1.0
	a := []float32{3.0, 4.0}
	b := []float32{6.0, 8.0}
	sim, err := vmath.CosineSimilarity(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(float64(sim)-1.0) > 1e-5 {
		t.Errorf("scaled vectors: expected 1.0, got %.6f", sim)
	}
}

func TestCosineSimilarity_LengthMismatch(t *testing.T) {
	a := []float32{1.0, 2.0}
	b := []float32{1.0, 2.0, 3.0}
	_, err := vmath.CosineSimilarity(a, b)
	if err == nil {
		t.Error("expected error for length mismatch, got nil")
	}
}

func TestCosineSimilarity_EmptyVectors(t *testing.T) {
	_, err := vmath.CosineSimilarity([]float32{}, []float32{})
	if err == nil {
		t.Error("expected error for empty vectors, got nil")
	}
}
