package main

import (
	"path/filepath"
	"testing"

	"github.com/zhangchuqi1998/velosearch/internal/collection"
	"github.com/zhangchuqi1998/velosearch/internal/storage"
)

func TestReplay_RestoresState(t *testing.T) {
	walPath := filepath.Join(t.TempDir(), "wal.log")

	// --- Session 1: write some records via WAL ---
	w, err := storage.Open(walPath)
	if err != nil {
		t.Fatalf("Open session 1: %v", err)
	}
	records := []*storage.Record{
		{Op: storage.OpCreateColl, CreateColl: &storage.CreateColl{
			Name: "c1", Dim: 4, Metric: 0, M: 8, EfConstruction: 50,
		}},
		{Op: storage.OpInsert, Insert: &storage.Insert{
			Collection: "c1", ID: 1, Vector: []float32{1, 2, 3, 4},
		}},
		{Op: storage.OpInsert, Insert: &storage.Insert{
			Collection: "c1", ID: 2, Vector: []float32{5, 6, 7, 8},
		}},
		{Op: storage.OpInsert, Insert: &storage.Insert{
			Collection: "c1", ID: 3, Vector: []float32{9, 10, 11, 12},
		}},
		{Op: storage.OpDelete, Delete: &storage.Delete{Collection: "c1", ID: 2}},
	}
	for i, r := range records {
		if err := w.Append(r); err != nil {
			t.Fatalf("Append %d: %v", i, err)
		}
	}
	w.Close()

	// --- Session 2: simulated restart — replay WAL into a fresh Manager ---
	mgr := collection.NewManager()
	n := 0
	if err := storage.Replay(walPath, func(rec *storage.Record) error {
		n++
		return ApplyWALRecord(mgr, rec)
	}); err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if n != len(records) {
		t.Errorf("replayed records: got %d, want %d", n, len(records))
	}

	col, err := mgr.Get("c1")
	if err != nil {
		t.Fatalf("Get c1 after replay: %v", err)
	}
	st := col.Index.SnapshotStats()
	if st.NumVectors != 3 {
		t.Errorf("NumVectors: got %d, want 3", st.NumVectors)
	}
	if st.NumDeleted != 1 {
		t.Errorf("NumDeleted: got %d, want 1", st.NumDeleted)
	}

	// Search for vector 1 — should find id=1, not id=2 (deleted).
	hits := col.Index.Search([]float32{1, 2, 3, 4}, 2, 10)
	if len(hits) == 0 {
		t.Fatal("Search: empty hits after replay")
	}
	if hits[0].ID != 1 {
		t.Errorf("Search top hit after replay: got id=%d, want 1", hits[0].ID)
	}
	for _, h := range hits {
		if h.ID == 2 {
			t.Errorf("deleted id=2 leaked into Search results")
		}
	}
}

func TestReplay_IdempotentOnRepeatedCreate(t *testing.T) {
	walPath := filepath.Join(t.TempDir(), "wal.log")

	w, _ := storage.Open(walPath)
	create := &storage.Record{Op: storage.OpCreateColl, CreateColl: &storage.CreateColl{
		Name: "c1", Dim: 4, Metric: 0, M: 8, EfConstruction: 50,
	}}
	if err := w.Append(create); err != nil {
		t.Fatalf("Append 1: %v", err)
	}
	if err := w.Append(create); err != nil {
		t.Fatalf("Append 2: %v", err)
	}
	w.Close()

	mgr := collection.NewManager()
	if err := storage.Replay(walPath, func(rec *storage.Record) error {
		return ApplyWALRecord(mgr, rec)
	}); err != nil {
		t.Fatalf("Replay: %v", err)
	}
	if _, err := mgr.Get("c1"); err != nil {
		t.Errorf("c1 not present after duplicate Create replay: %v", err)
	}
}
