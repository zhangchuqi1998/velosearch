package sift

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

// TestLoadFvecs_SIFTBase asserts the canonical SIFT-1M base file loads with
// the expected shape. Skipped when the dataset isn't present so CI passes
// without downloading 170 MB.
func TestLoadFvecs_SIFTBase(t *testing.T) {
	path := "data/sift_base.fvecs"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("SIFT base file not found at %s — skip (download per Day 5 roadmap)", path)
	}
	got, err := LoadFvecs(path)
	if err != nil {
		t.Fatalf("LoadFvecs: %v", err)
	}
	if want := 1_000_000; len(got) != want {
		t.Errorf("base vector count = %d, want %d", len(got), want)
	}
	if want := 128; len(got[0]) != want {
		t.Errorf("base vector dim = %d, want %d", len(got[0]), want)
	}
}

// TestLoadFvecs_SIFTQuery asserts the query file shape.
func TestLoadFvecs_SIFTQuery(t *testing.T) {
	path := "data/sift_query.fvecs"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("SIFT query file not found at %s", path)
	}
	got, err := LoadFvecs(path)
	if err != nil {
		t.Fatalf("LoadFvecs: %v", err)
	}
	if want := 10_000; len(got) != want {
		t.Errorf("query vector count = %d, want %d", len(got), want)
	}
	if want := 128; len(got[0]) != want {
		t.Errorf("query vector dim = %d, want %d", len(got[0]), want)
	}
}

// TestLoadIvecs_SIFTGroundTruth asserts the ground-truth neighbor file shape.
func TestLoadIvecs_SIFTGroundTruth(t *testing.T) {
	path := "data/sift_groundtruth.ivecs"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("SIFT ground-truth file not found at %s", path)
	}
	got, err := LoadIvecs(path)
	if err != nil {
		t.Fatalf("LoadIvecs: %v", err)
	}
	if want := 10_000; len(got) != want {
		t.Errorf("gt vector count = %d, want %d", len(got), want)
	}
	if want := 100; len(got[0]) != want {
		t.Errorf("gt neighbors per query = %d, want %d", len(got[0]), want)
	}
}

// TestLoadFvecs_Synthetic builds a tiny .fvecs file by hand and round-trips
// it through LoadFvecs. Runs in CI without any external data.
func TestLoadFvecs_Synthetic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tiny.fvecs")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]float32{
		{1, 2, 3, 4},
		{-5, 0.5, 99, -0.25},
		{0, 0, 0, 0},
	}
	for _, v := range want {
		if err := binary.Write(f, binary.LittleEndian, int32(len(v))); err != nil {
			t.Fatal(err)
		}
		if err := binary.Write(f, binary.LittleEndian, v); err != nil {
			t.Fatal(err)
		}
	}
	f.Close()

	got, err := LoadFvecs(path)
	if err != nil {
		t.Fatalf("LoadFvecs: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("vector count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if len(got[i]) != len(want[i]) {
			t.Fatalf("vec %d: got len %d, want %d", i, len(got[i]), len(want[i]))
		}
		for j := range want[i] {
			if got[i][j] != want[i][j] {
				t.Errorf("vec %d[%d] = %v, want %v", i, j, got[i][j], want[i][j])
			}
		}
	}
}

// TestLoadIvecs_Synthetic same as above for ivecs.
func TestLoadIvecs_Synthetic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tiny.ivecs")

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	want := [][]int32{
		{10, 20, 30},
		{-1, 0, 1},
	}
	for _, v := range want {
		if err := binary.Write(f, binary.LittleEndian, int32(len(v))); err != nil {
			t.Fatal(err)
		}
		if err := binary.Write(f, binary.LittleEndian, v); err != nil {
			t.Fatal(err)
		}
	}
	f.Close()

	got, err := LoadIvecs(path)
	if err != nil {
		t.Fatalf("LoadIvecs: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("vector count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		for j := range want[i] {
			if got[i][j] != want[i][j] {
				t.Errorf("vec %d[%d] = %v, want %v", i, j, got[i][j], want[i][j])
			}
		}
	}
}

// TestLoadFvecs_FileNotFound checks error path.
func TestLoadFvecs_FileNotFound(t *testing.T) {
	if _, err := LoadFvecs("/nonexistent/path/foo.fvecs"); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
