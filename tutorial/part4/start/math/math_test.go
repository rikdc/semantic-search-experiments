package math

import (
	"math"
	"testing"
)

// --- FormatTransaction ---

func TestFormatTransaction(t *testing.T) {
	tx := Transaction{
		Label: "Normal_1", Category: "Normal",
		Amount: 45.00, Merchant: "Whole Foods Market", Time: "11:00 AM",
	}
	got := FormatTransaction(tx)
	want := "Transaction: 45.00 USD at Whole Foods Market at 11:00 AM"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatTransaction_HighValueLateNight(t *testing.T) {
	tx := Transaction{
		Label: "Anomaly", Category: "Fraud_Candidate",
		Amount: 4999.00, Merchant: "High-End Jewelry Store", Time: "03:45 AM",
	}
	got := FormatTransaction(tx)
	want := "Transaction: 4999.00 USD at High-End Jewelry Store at 03:45 AM"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// --- FormatRaw ---

func TestFormatRaw(t *testing.T) {
	tx := Transaction{Amount: 4999.00}
	got := FormatRaw(tx)
	want := "4999.00"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatRaw_SmallAmount(t *testing.T) {
	tx := Transaction{Amount: 8.75}
	got := FormatRaw(tx)
	want := "8.75"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// --- Detection pipeline with hand-crafted vectors ---
//
// These tests verify the centroid + threshold logic before any API call.
// The vectors are chosen so the correct answer is obvious: normals point
// in one direction, the outlier points in roughly the opposite direction.

func TestDetection_NormalsPassThreshold(t *testing.T) {
	normalA := []float32{0.9, 0.1, 0.8, 0.2}
	normalB := []float32{0.85, 0.15, 0.75, 0.25}
	normalC := []float32{0.88, 0.12, 0.82, 0.18}
	outlier := []float32{0.1, 0.9, 0.2, 0.8}

	centroid, err := CalculateCentroid([][]float32{normalA, normalB, normalC})
	if err != nil {
		t.Fatalf("centroid: %v", err)
	}

	const threshold = float32(0.80)

	for _, v := range [][]float32{normalA, normalB, normalC} {
		sim, err := CosineSimilarity(v, centroid)
		if err != nil {
			t.Fatal(err)
		}
		if sim < threshold {
			t.Errorf("normal vector incorrectly flagged: sim=%.4f, threshold=%.2f", sim, threshold)
		}
	}

	simO, err := CosineSimilarity(outlier, centroid)
	if err != nil {
		t.Fatal(err)
	}
	if simO >= threshold {
		t.Errorf("outlier not flagged: sim=%.4f, threshold=%.2f", simO, threshold)
	}
}

func TestDetection_SelfSimilarityIsOne(t *testing.T) {
	// A vector compared to itself must be exactly 1.0.
	// float64 accumulation in CosineSimilarity ensures this;
	// float32 accumulation would drift on 1536-dimensional vectors.
	v := []float32{0.9, 0.1, 0.8, 0.2}
	sim, err := CosineSimilarity(v, v)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(float64(sim)-1.0) > 1e-6 {
		t.Errorf("self-similarity = %.8f, want 1.0", sim)
	}
}

func TestDetection_CentroidFromOneVector(t *testing.T) {
	// Centroid of a single vector is that vector.
	v := []float32{0.5, 0.5, 0.5, 0.5}
	c, err := CalculateCentroid([][]float32{v})
	if err != nil {
		t.Fatal(err)
	}
	for i := range v {
		if math.Abs(float64(c[i]-v[i])) > 1e-6 {
			t.Errorf("centroid[%d] = %f, want %f", i, c[i], v[i])
		}
	}
}
