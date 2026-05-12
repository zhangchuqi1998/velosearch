// Package distance provides vector distance and similarity functions
// commonly used in approximate nearest neighbor (ANN) search.
package distance

import "math"

// DistanceFunc is the signature for any distance/similarity function
// operating on two equal-length float32 vectors.
type DistanceFunc func(a, b []float32) float32

// L2Squared returns the squared Euclidean distance between a and b.
// The square root is intentionally omitted: for nearest-neighbor
// ranking, squared distance preserves order and is cheaper to compute.
// Panics if len(a) != len(b).
func L2Squared(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("distance: L2Squared: length mismatch")
	}
	var sum float32
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return sum
}

// DotProduct returns the dot product of a and b.
// Panics if len(a) != len(b).
func DotProduct(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("distance: DotProduct: length mismatch")
	}
	var sum float32
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// Norm returns the Euclidean (L2) norm of a, i.e. sqrt(sum(a[i]^2)).
func Norm(a []float32) float32 {
	var sum float32
	for _, v := range a {
		sum += v * v
	}
	return float32(math.Sqrt(float64(sum)))
}

// Cosine returns the cosine distance between a and b, defined as
// 1 - cos(a, b) = 1 - dot(a, b) / (|a| * |b|).
// Range is [0, 2]: identical vectors return 0, orthogonal return 1,
// opposite return 2. If either vector is zero, returns 1 (treating
// undefined similarity as maximally dissimilar within the orthogonal case).
// Panics if len(a) != len(b).
func Cosine(a, b []float32) float32 {
	if len(a) != len(b) {
		panic("distance: Cosine: length mismatch")
	}
	var dot, na, nb float32
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	denom := float32(math.Sqrt(float64(na))) * float32(math.Sqrt(float64(nb)))
	if denom == 0 {
		return 1
	}
	return 1 - dot/denom
}
