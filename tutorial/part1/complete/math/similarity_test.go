package math

import (
	"math"
	"testing"
)

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name      string
		a         []float32
		b         []float32
		want      float32
		wantErr   bool
		tolerance float32
	}{
		{
			name: "identical vectors",
			a:    []float32{1, 2, 3},
			b:    []float32{1, 2, 3},
			want: 1.0,
		},
		{
			name: "opposite vectors",
			a:    []float32{1, 2, 3},
			b:    []float32{-1, -2, -3},
			want: -1.0,
		},
		{
			name: "orthogonal vectors",
			a:    []float32{1, 0},
			b:    []float32{0, 1},
			want: 0.0,
		},
		{
			name:      "proportional vectors",
			a:         []float32{1, 2, 3},
			b:         []float32{2, 4, 6},
			want:      1.0,
			tolerance: 0.0001,
		},
		{
			name:      "somewhat similar",
			a:         []float32{1, 2, 3},
			b:         []float32{1, 2, 1},
			want:      0.8729,
			tolerance: 0.001,
		},
		{
			name:    "mismatched dimensions",
			a:       []float32{1, 2, 3},
			b:       []float32{1, 2},
			wantErr: true,
		},
		{
			name:    "zero vector first",
			a:       []float32{0, 0, 0},
			b:       []float32{1, 2, 3},
			wantErr: true,
		},
		{
			name:    "zero vector second",
			a:       []float32{1, 2, 3},
			b:       []float32{0, 0, 0},
			wantErr: true,
		},
		{
			name:    "both zero vectors",
			a:       []float32{0, 0, 0},
			b:       []float32{0, 0, 0},
			wantErr: true,
		},
		{
			name:    "empty vectors",
			a:       []float32{},
			b:       []float32{},
			wantErr: true,
		},
		{
			name:      "single dimension",
			a:         []float32{5},
			b:         []float32{10},
			want:      1.0,
			tolerance: 0.0001,
		},
		{
			name:      "negative values",
			a:         []float32{-1, -2, -3},
			b:         []float32{-2, -4, -6},
			want:      1.0,
			tolerance: 0.0001,
		},
		{
			name:      "mixed positive and negative",
			a:         []float32{1, -1, 2},
			b:         []float32{-1, 1, -2},
			want:      -1.0,
			tolerance: 0.0001,
		},
		{
			name:      "very small values (precision test)",
			a:         []float32{0.0001, 0.0002, 0.0003},
			b:         []float32{0.0002, 0.0004, 0.0006},
			want:      1.0,
			tolerance: 0.0001,
		},
		{
			name:      "product vectors from walkthrough - similar",
			a:         []float32{8.5, 1.2, 9.2},
			b:         []float32{8.3, 1.5, 9.0},
			want:      0.9997,
			tolerance: 0.001,
		},
		{
			name:      "product vectors from walkthrough - different",
			a:         []float32{8.5, 1.2, 9.2},
			b:         []float32{1.0, 9.5, 0.5},
			want:      0.2036,
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CosineSimilarity(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("CosineSimilarity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				tolerance := tt.tolerance
				if tolerance == 0 {
					tolerance = 0.0001
				}
				if math.Abs(float64(got-tt.want)) > float64(tolerance) {
					t.Errorf("CosineSimilarity() = %v, want %v (tolerance: %v)", got, tt.want, tolerance)
				}
			}
		})
	}
}

func TestCosineSimilarityHighDimensional(t *testing.T) {
	a := make([]float32, 1536)
	b := make([]float32, 1536)
	for i := range a {
		a[i] = float32(i)
		b[i] = float32(i)
	}
	// Skip index 0 which is zero in both — that's fine, non-zero elements dominate.
	got, err := CosineSimilarity(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(float64(got-1.0)) > 0.0001 {
		t.Errorf("identical 1536-dim vectors should have similarity 1.0, got %v", got)
	}
}

func TestCosineSimilarityCommonMistakes(t *testing.T) {
	t.Run("proportional vectors score 1.0", func(t *testing.T) {
		a := []float32{3, 4}
		b := []float32{6, 8}
		got, err := CosineSimilarity(a, b)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if math.Abs(float64(got-1.0)) > 0.0001 {
			t.Errorf("expected 1.0 for proportional vectors, got %v", got)
		}
	})

	t.Run("does not modify input vectors", func(t *testing.T) {
		a := []float32{1, 2, 3}
		b := []float32{4, 5, 6}
		aCopy := make([]float32, len(a))
		bCopy := make([]float32, len(b))
		copy(aCopy, a)
		copy(bCopy, b)

		_, err := CosineSimilarity(a, b)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		for i := range a {
			if a[i] != aCopy[i] {
				t.Errorf("input vector a modified at index %d: %v != %v", i, a[i], aCopy[i])
			}
		}
		for i := range b {
			if b[i] != bCopy[i] {
				t.Errorf("input vector b modified at index %d: %v != %v", i, b[i], bCopy[i])
			}
		}
	})

	t.Run("float32 precision with identical vectors", func(t *testing.T) {
		a := []float32{0.1, 0.2, 0.3, 0.4}
		b := []float32{0.1, 0.2, 0.3, 0.4}
		got, err := CosineSimilarity(a, b)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if math.Abs(float64(got-1.0)) > 0.001 {
			t.Errorf("identical vectors should have similarity 1.0, got %v (likely using float32 for intermediates)", got)
		}
	})
}

func TestCalculateCentroid(t *testing.T) {
	tests := []struct {
		name    string
		vectors [][]float32
		want    []float32
		wantErr bool
	}{
		{
			name: "single vector",
			vectors: [][]float32{
				{1, 2, 3},
			},
			want: []float32{1, 2, 3},
		},
		{
			name: "two vectors",
			vectors: [][]float32{
				{2, 4},
				{4, 6},
			},
			want: []float32{3, 5},
		},
		{
			name: "three vectors",
			vectors: [][]float32{
				{2, 4},
				{4, 6},
				{6, 8},
			},
			want: []float32{4, 6},
		},
		{
			name:    "empty input",
			vectors: [][]float32{},
			wantErr: true,
		},
		{
			name: "mismatched dimensions",
			vectors: [][]float32{
				{1, 2, 3},
				{4, 5},
			},
			wantErr: true,
		},
		{
			name: "negative values",
			vectors: [][]float32{
				{-1, -2},
				{1, 2},
			},
			want: []float32{0, 0},
		},
		{
			name: "identical vectors",
			vectors: [][]float32{
				{5, 10, 15},
				{5, 10, 15},
				{5, 10, 15},
			},
			want: []float32{5, 10, 15},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateCentroid(tt.vectors)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateCentroid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Fatalf("CalculateCentroid() returned %d dims, want %d", len(got), len(tt.want))
				}
				for i := range got {
					if math.Abs(float64(got[i]-tt.want[i])) > 0.0001 {
						t.Errorf("CalculateCentroid()[%d] = %v, want %v", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}
