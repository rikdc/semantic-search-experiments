package math

import (
	"math"
	"testing"
)

// --- CalculateMeanStdDev ---

func TestCalculateMeanStdDev_KnownValues(t *testing.T) {
	// mean = (0.80 + 0.84 + 0.86 + 0.90) / 4 = 0.85
	// stddev (population) ≈ 0.036
	sims := []float32{0.80, 0.84, 0.86, 0.90}
	mean, stddev := CalculateMeanStdDev(sims)

	if math.Abs(mean-0.85) > 1e-5 {
		t.Errorf("mean = %.6f, want ≈ 0.85", mean)
	}
	if math.Abs(stddev-0.036) > 1e-3 {
		t.Errorf("stddev = %.6f, want ≈ 0.036", stddev)
	}
}

func TestCalculateMeanStdDev_IdenticalValues(t *testing.T) {
	// All identical inputs: mean = that value, stddev = 0
	sims := []float32{0.90, 0.90, 0.90}
	mean, stddev := CalculateMeanStdDev(sims)

	if math.Abs(mean-0.90) > 1e-6 {
		t.Errorf("mean = %.6f, want 0.90", mean)
	}
	if stddev != 0 {
		t.Errorf("stddev = %.6f, want 0 for identical inputs", stddev)
	}
}

func TestCalculateMeanStdDev_SingleValue(t *testing.T) {
	sims := []float32{0.85}
	mean, stddev := CalculateMeanStdDev(sims)

	if math.Abs(mean-0.85) > 1e-6 {
		t.Errorf("mean = %.6f, want 0.85", mean)
	}
	if stddev != 0 {
		t.Errorf("stddev = %.6f, want 0 for single value", stddev)
	}
}

func TestCalculateMeanStdDev_EmptySlice(t *testing.T) {
	mean, stddev := CalculateMeanStdDev([]float32{})
	if mean != 0 || stddev != 0 {
		t.Errorf("empty slice: got mean=%.4f stddev=%.4f, want both 0", mean, stddev)
	}
}

// --- Z-score detection with hand-crafted similarities ---
//
// These tests verify the full z-score flagging pipeline before any API call.
// The similarities are chosen so the expected results are obvious by inspection.

func TestZScore_OutlierFlagged(t *testing.T) {
	// Normal cluster: tight group around 0.85
	// Outlier: 0.60 — clearly far from the cluster
	normalSims := []float32{0.80, 0.84, 0.86, 0.90}
	outlierSim := float32(0.60)

	mean, stddev := CalculateMeanStdDev(normalSims)

	zOutlier := (float64(outlierSim) - mean) / stddev
	const threshold = -2.0

	if zOutlier >= threshold {
		t.Errorf("outlier not flagged: z=%.2f, threshold=%.1f", zOutlier, threshold)
	}
}

func TestZScore_NormalsNotFlagged(t *testing.T) {
	normalSims := []float32{0.80, 0.84, 0.86, 0.90}
	mean, stddev := CalculateMeanStdDev(normalSims)

	const threshold = -2.0
	for _, s := range normalSims {
		z := (float64(s) - mean) / stddev
		if z < threshold {
			t.Errorf("normal sim %.4f incorrectly flagged: z=%.2f", s, z)
		}
	}
}

// --- Inherited from Part 4: format and centroid sanity checks ---

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

func TestFormatRaw(t *testing.T) {
	tx := Transaction{Amount: 4999.00}
	got := FormatRaw(tx)
	want := "4999.00"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDetection_SelfSimilarityIsOne(t *testing.T) {
	v := []float32{0.9, 0.1, 0.8, 0.2}
	sim, err := CosineSimilarity(v, v)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(float64(sim)-1.0) > 1e-6 {
		t.Errorf("self-similarity = %.8f, want 1.0", sim)
	}
}
