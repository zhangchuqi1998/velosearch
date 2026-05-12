// Package distance provides vector distance and similarity functions
// commonly used in approximate nearest neighbor (ANN) search.
//
// The implementations dispatch to AVX2-accelerated routines from
// github.com/viterin/vek/vek32 when the CPU supports it, and fall back
// to a portable Go loop otherwise (transparently handled by vek).
package distance

import (
	"math"

	"github.com/viterin/vek/vek32"
)

// DistanceFunc is the signature for any distance/similarity function
// operating on two equal-length float32 vectors.
type DistanceFunc func(a, b []float32) float32

// L2Squared returns the squared Euclidean distance between a and b.
// For nearest-neighbor ranking, squared distance preserves order and
// avoids an unnecessary sqrt at the call site.
//
// Internally we call vek32.Distance (SIMD, returns sqrt'd distance)
// and re-square — the SIMD subtract+square+sum dominates the cost,
// and the redundant sqrt+square is negligible (~1 ns).
//
// Panics if len(a) != len(b).
func L2Squared(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("distance: L2Squared: length mismatch")
	}
	if len(a) == 0 {
		return 0
	}
	d := vek32.Distance(a, b)
	return d * d
}

// DotProduct returns the dot product of a and b. Panics if lengths differ.
func DotProduct(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("distance: DotProduct: length mismatch")
	}
	if len(a) == 0 {
		return 0
	}
	return vek32.Dot(a, b)
}

// Norm returns the Euclidean (L2) norm of a, i.e. sqrt(sum(a[i]^2)).
func Norm(a []float32) float32 {
	if len(a) == 0 {
		return 0
	}
	// vek32 has no direct Norm; use Dot(a, a) then sqrt.
	return float32(math.Sqrt(float64(vek32.Dot(a, a))))
}

// Cosine returns the cosine distance between a and b, defined as
// 1 - cos(a, b) = 1 - dot(a, b) / (|a| * |b|).
// Range is [0, 2]: identical vectors return 0, orthogonal return 1,
// opposite return 2. If either vector is zero, returns 1.
// Panics if len(a) != len(b).
func Cosine(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("distance: Cosine: length mismatch")
	}
	if len(a) == 0 {
		return 1
	}
	// vek32.CosineSimilarity returns cos(a, b) in [-1, 1].
	// Defined as 0 when either operand has zero norm — we mirror our
	// previous "return 1" convention here for explicitness.
	na := vek32.Dot(a, a)
	nb := vek32.Dot(b, b)
	if na == 0 || nb == 0 {
		return 1
	}
	return 1 - vek32.CosineSimilarity(a, b)
}
