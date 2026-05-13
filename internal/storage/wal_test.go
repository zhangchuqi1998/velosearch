package storage

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func sampleRecords() []*Record {
	return []*Record{
		{Op: OpCreateColl, CreateColl: &CreateColl{Name: "c1", Dim: 8, Metric: 0, M: 16, EfConstruction: 200}},
		{Op: OpInsert, Insert: &Insert{Collection: "c1", ID: 1, Vector: []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}}},
		{Op: OpInsert, Insert: &Insert{Collection: "c1", ID: 2, Vector: []float32{1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8}}},
		{Op: OpDelete, Delete: &Delete{Collection: "c1", ID: 1}},
		{Op: OpDropColl, DropColl: &DropColl{Name: "c1"}},
	}
}

func TestWAL_AppendReplay(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wal.log")

	w, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	want := sampleRecords()
	for i, rec := range want {
		if err := w.Append(rec); err != nil {
			t.Fatalf("Append %d: %v", i, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	var got []*Record
	if err := Replay(path, func(r *Record) error {
		got = append(got, r)
		return nil
	}); err != nil {
		t.Fatalf("Replay: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("record count: got %d, want %d", len(got), len(want))
	}
	for i := range want {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Errorf("record %d:\n got=%+v\nwant=%+v", i, got[i], want[i])
		}
	}
}

func TestWAL_TruncatedTail(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wal.log")

	w, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	recs := sampleRecords()
	for _, r := range recs {
		if err := w.Append(r); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}
	w.Close()

	stat, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if err := os.Truncate(path, stat.Size()-10); err != nil {
		t.Fatalf("Truncate: %v", err)
	}

	var n int
	err = Replay(path, func(r *Record) error {
		n++
		return nil
	})
	if err != nil {
		t.Fatalf("Replay should warn-and-stop on truncated tail, got error: %v", err)
	}
	if n != len(recs)-1 {
		t.Errorf("record count after truncate: got %d, want %d", n, len(recs)-1)
	}
}

func TestWAL_CorruptedBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wal.log")

	w, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	recs := sampleRecords()
	for _, r := range recs {
		if err := w.Append(r); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}
	w.Close()

	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	var lenBuf [4]byte
	if _, err := f.ReadAt(lenBuf[:], 0); err != nil {
		t.Fatalf("ReadAt len: %v", err)
	}
	len0 := binary.BigEndian.Uint32(lenBuf[:])
	flipAt := int64(4) + int64(len0)/2
	var b [1]byte
	if _, err := f.ReadAt(b[:], flipAt); err != nil {
		t.Fatalf("ReadAt: %v", err)
	}
	b[0] ^= 0xFF
	if _, err := f.WriteAt(b[:], flipAt); err != nil {
		t.Fatalf("WriteAt: %v", err)
	}
	f.Close()

	err = Replay(path, func(r *Record) error { return nil })
	if err == nil {
		t.Fatal("Replay should fail on CRC mismatch, got nil")
	}
	if !strings.Contains(err.Error(), "crc") {
		t.Errorf("expected crc error, got: %v", err)
	}
}
