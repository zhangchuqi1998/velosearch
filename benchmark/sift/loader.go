// Package sift loads the SIFT-1M vector benchmark dataset distributed by
// http://corpus-texmex.irisa.fr/.
//
// The dataset uses two binary formats:
//
//   .fvecs : sequence of float32 vectors
//   .ivecs : sequence of int32 vectors (used for ground-truth neighbor IDs)
//
// Each vector in both formats is preceded by a 4-byte little-endian int32
// holding the dimensionality d, followed by d × 4 bytes of values. Vectors
// are concatenated until EOF.
//
// Typical files in the SIFT-1M corpus:
//
//   sift_base.fvecs         1,000,000 × 128 float32
//   sift_query.fvecs           10,000 × 128 float32
//   sift_groundtruth.ivecs     10,000 × 100 int32
//   sift_learn.fvecs          100,000 × 128 float32 (not used)
package sift

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// LoadFvecs loads a .fvecs file into a [][]float32. Every vector is expected
// to have the same dimensionality; an error is returned otherwise.
func LoadFvecs(path string) ([][]float32, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	br := bufio.NewReaderSize(f, 1<<20) // 1 MB buffer

	var (
		out    [][]float32
		dimRef int32 = -1
	)
	for {
		var dim int32
		if err := binary.Read(br, binary.LittleEndian, &dim); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return nil, fmt.Errorf("read dim: %w", err)
		}
		if dim <= 0 {
			return nil, fmt.Errorf("non-positive dim %d at vector %d", dim, len(out))
		}
		if dimRef < 0 {
			dimRef = dim
		} else if dim != dimRef {
			return nil, fmt.Errorf("inconsistent dim at vector %d: got %d, want %d",
				len(out), dim, dimRef)
		}

		v := make([]float32, dim)
		if err := binary.Read(br, binary.LittleEndian, v); err != nil {
			return nil, fmt.Errorf("read vector %d: %w", len(out), err)
		}
		out = append(out, v)
	}
	return out, nil
}

// LoadIvecs loads a .ivecs file into a [][]int32.
func LoadIvecs(path string) ([][]int32, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	br := bufio.NewReaderSize(f, 1<<20)

	var (
		out    [][]int32
		dimRef int32 = -1
	)
	for {
		var dim int32
		if err := binary.Read(br, binary.LittleEndian, &dim); err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			return nil, fmt.Errorf("read dim: %w", err)
		}
		if dim <= 0 {
			return nil, fmt.Errorf("non-positive dim %d at vector %d", dim, len(out))
		}
		if dimRef < 0 {
			dimRef = dim
		} else if dim != dimRef {
			return nil, fmt.Errorf("inconsistent dim at vector %d: got %d, want %d",
				len(out), dim, dimRef)
		}

		v := make([]int32, dim)
		if err := binary.Read(br, binary.LittleEndian, v); err != nil {
			return nil, fmt.Errorf("read vector %d: %w", len(out), err)
		}
		out = append(out, v)
	}
	return out, nil
}
