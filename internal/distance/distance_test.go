package distance

import (
	"math"
	"math/rand"
	"testing"
)

const epsilon = 1e-6

func floatsClose(a, b, tol float32) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tol
}

func TestL2Squared(t *testing.T) {
	tests := []struct {
		name string
		a, b []float32
		want float32
	}{
		{"3-4-5 triangle", []float32{0, 0}, []float32{3, 4}, 25},
		{"identical", []float32{1, 2, 3}, []float32{1, 2, 3}, 0},
		{"unit x-axis", []float32{0, 0, 0}, []float32{1, 0, 0}, 1},
		{"negatives", []float32{-1, -1}, []float32{1, 1}, 8},
		{"empty", []float32{}, []float32{}, 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := L2Squared(tc.a, tc.b)
			if !floatsClose(got, tc.want, epsilon) {
				t.Errorf("L2Squared(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestCosine(t *testing.T) {
	tests := []struct {
		name string
		a, b []float32
		want float32
	}{
		{"orthogonal", []float32{1, 0}, []float32{0, 1}, 1.0},
		{"identical", []float32{1, 2, 3}, []float32{1, 2, 3}, 0.0},
		{"identical unit", []float32{1, 0, 0}, []float32{1, 0, 0}, 0.0},
		{"opposite", []float32{1, 0}, []float32{-1, 0}, 2.0},
		{"scaled identical", []float32{1, 2, 3}, []float32{2, 4, 6}, 0.0},
		{"zero vector", []float32{0, 0, 0}, []float32{1, 2, 3}, 1.0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Cosine(tc.a, tc.b)
			if !floatsClose(got, tc.want, epsilon) {
				t.Errorf("Cosine(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestDotProduct(t *testing.T) {
	tests := []struct {
		name string
		a, b []float32
		want float32
	}{
		{"basic", []float32{1, 2, 3}, []float32{4, 5, 6}, 32},
		{"orthogonal", []float32{1, 0}, []float32{0, 1}, 0},
		{"with negatives", []float32{1, -2}, []float32{-3, 4}, -11},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := DotProduct(tc.a, tc.b)
			if !floatsClose(got, tc.want, epsilon) {
				t.Errorf("DotProduct(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestNorm(t *testing.T) {
	tests := []struct {
		name string
		a    []float32
		want float32
	}{
		{"3-4-5", []float32{3, 4}, 5},
		{"unit", []float32{1, 0, 0}, 1},
		{"zero", []float32{0, 0, 0}, 0},
		{"sqrt(14)", []float32{1, 2, 3}, float32(math.Sqrt(14))},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Norm(tc.a)
			if !floatsClose(got, tc.want, epsilon) {
				t.Errorf("Norm(%v) = %v, want %v", tc.a, got, tc.want)
			}
		})
	}
}

func TestPanicOnLengthMismatch(t *testing.T) {
	cases := []struct {
		name string
		fn   func()
	}{
		{"L2Squared", func() { L2Squared([]float32{1, 2}, []float32{1, 2, 3}) }},
		{"Cosine", func() { Cosine([]float32{1, 2}, []float32{1, 2, 3}) }},
		{"DotProduct", func() { DotProduct([]float32{1, 2}, []float32{1, 2, 3}) }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("%s: expected panic on length mismatch, got none", tc.name)
				}
			}()
			tc.fn()
		})
	}
}

// TestDistanceFuncType verifies the function signatures conform to DistanceFunc.
func TestDistanceFuncType(t *testing.T) {
	var fns = []DistanceFunc{L2Squared, Cosine, DotProduct}
	a := []float32{1, 2, 3}
	b := []float32{4, 5, 6}
	for _, fn := range fns {
		_ = fn(a, b) // smoke check: should not panic
	}
}

// --- Benchmarks ---

const benchDim = 128

func makeVec(seed int64, n int) []float32 {
	r := rand.New(rand.NewSource(seed))
	v := make([]float32, n)
	for i := range v {
		v[i] = r.Float32()
	}
	return v
}

var (
	benchA = makeVec(1, benchDim)
	benchB = makeVec(2, benchDim)
	// Sinks to prevent the compiler from optimizing benchmark calls away.
	sinkF32 float32
)

func BenchmarkL2Squared(b *testing.B) {
	b.ReportAllocs()
	var s float32
	for i := 0; i < b.N; i++ {
		s = L2Squared(benchA, benchB)
	}
	sinkF32 = s
}

func BenchmarkCosine(b *testing.B) {
	b.ReportAllocs()
	var s float32
	for i := 0; i < b.N; i++ {
		s = Cosine(benchA, benchB)
	}
	sinkF32 = s
}
