// Package storage implements VeloSearch's write-ahead log.
//
// On-disk format:
//
//	+----------+----------+----------+
//	| 4 bytes  | N bytes  | 4 bytes  |
//	| BE len N | payload  | CRC32 BE |
//	+----------+----------+----------+
//	... repeated ...
//
// Payload layout:
//
//	[1 byte] op type (1=CreateColl, 2=DropColl, 3=Insert, 4=Delete)
//	[op-specific bytes]
//
// CRC32 uses the Castagnoli polynomial (same as ext4, btrfs, iSCSI).
package storage

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log/slog"
	"math"
	"os"
	"sync"
)

type OpType uint8

const (
	OpCreateColl OpType = 1
	OpDropColl   OpType = 2
	OpInsert     OpType = 3
	OpDelete     OpType = 4
)

type Record struct {
	Op         OpType
	CreateColl *CreateColl
	DropColl   *DropColl
	Insert     *Insert
	Delete     *Delete
}

type CreateColl struct {
	Name           string
	Dim            int32
	Metric         uint8
	M              int32
	EfConstruction int32
}

type DropColl struct {
	Name string
}

type Insert struct {
	Collection string
	ID         uint32
	Vector     []float32
}

type Delete struct {
	Collection string
	ID         uint32
}

var crcTable = crc32.MakeTable(crc32.Castagnoli)

type WAL struct {
	mu sync.Mutex
	f  *os.File
}

func Open(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open wal %s: %w", path, err)
	}
	return &WAL{f: f}, nil
}

// Append is safe to call on a nil receiver — it becomes a no-op. This lets
// callers run with WAL disabled (e.g. benchmarking) without branching at
// every call site.
func (w *WAL) Append(rec *Record) error {
	if w == nil {
		return nil
	}
	payload, err := marshal(rec)
	if err != nil {
		return err
	}
	crc := crc32.Checksum(payload, crcTable)

	frame := make([]byte, 4+len(payload)+4)
	binary.BigEndian.PutUint32(frame[:4], uint32(len(payload)))
	copy(frame[4:4+len(payload)], payload)
	binary.BigEndian.PutUint32(frame[4+len(payload):], crc)

	w.mu.Lock()
	defer w.mu.Unlock()
	if _, err := w.f.Write(frame); err != nil {
		return fmt.Errorf("wal write: %w", err)
	}
	if err := w.f.Sync(); err != nil {
		return fmt.Errorf("wal fsync: %w", err)
	}
	return nil
}

func (w *WAL) Close() error {
	if w == nil {
		return nil
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.f.Close()
}

// Replay reads records sequentially from path and invokes handler for each.
// A truncated tail (from a crash mid-write) logs a warning and returns nil.
// A CRC mismatch mid-file returns an error.
func Replay(path string, handler func(*Record) error) error {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("open wal %s: %w", path, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}
	totalSize := stat.Size()
	var offset int64

	for {
		var lenBuf [4]byte
		n, err := io.ReadFull(f, lenBuf[:])
		if err == io.EOF {
			return nil
		}
		if err == io.ErrUnexpectedEOF || (err == nil && n < 4) {
			slog.Warn("wal: truncated length prefix at end", "offset", offset, "have", n)
			return nil
		}
		if err != nil {
			return fmt.Errorf("wal read len at %d: %w", offset, err)
		}
		payloadLen := binary.BigEndian.Uint32(lenBuf[:])

		// Sanity-check: a single record larger than the remaining file is
		// either corruption or a truncated tail — treat as truncated and stop.
		remaining := totalSize - offset - 4
		if int64(payloadLen)+4 > remaining {
			slog.Warn("wal: truncated record",
				"offset", offset, "declared_len", payloadLen, "remaining", remaining)
			return nil
		}

		payload := make([]byte, payloadLen)
		if _, err := io.ReadFull(f, payload); err != nil {
			slog.Warn("wal: truncated payload", "offset", offset, "err", err)
			return nil
		}

		var crcBuf [4]byte
		if _, err := io.ReadFull(f, crcBuf[:]); err != nil {
			slog.Warn("wal: truncated crc", "offset", offset, "err", err)
			return nil
		}

		gotCRC := binary.BigEndian.Uint32(crcBuf[:])
		wantCRC := crc32.Checksum(payload, crcTable)
		if gotCRC != wantCRC {
			return fmt.Errorf("wal crc mismatch at offset %d: got %x want %x",
				offset, gotCRC, wantCRC)
		}

		rec, err := unmarshal(payload)
		if err != nil {
			return fmt.Errorf("wal unmarshal at offset %d: %w", offset, err)
		}

		if err := handler(rec); err != nil {
			return fmt.Errorf("wal handler at offset %d: %w", offset, err)
		}

		offset += int64(4 + payloadLen + 4)
	}
}

// --- marshal / unmarshal -----------------------------------------------------

func marshal(rec *Record) ([]byte, error) {
	switch rec.Op {
	case OpCreateColl:
		return marshalCreateColl(rec.CreateColl)
	case OpDropColl:
		return marshalDropColl(rec.DropColl)
	case OpInsert:
		return marshalInsert(rec.Insert)
	case OpDelete:
		return marshalDelete(rec.Delete)
	default:
		return nil, fmt.Errorf("unknown op type %d", rec.Op)
	}
}

func unmarshal(payload []byte) (*Record, error) {
	if len(payload) < 1 {
		return nil, errors.New("payload too short")
	}
	op := OpType(payload[0])
	body := payload[1:]
	switch op {
	case OpCreateColl:
		cc, err := unmarshalCreateColl(body)
		if err != nil {
			return nil, err
		}
		return &Record{Op: op, CreateColl: cc}, nil
	case OpDropColl:
		dc, err := unmarshalDropColl(body)
		if err != nil {
			return nil, err
		}
		return &Record{Op: op, DropColl: dc}, nil
	case OpInsert:
		ins, err := unmarshalInsert(body)
		if err != nil {
			return nil, err
		}
		return &Record{Op: op, Insert: ins}, nil
	case OpDelete:
		del, err := unmarshalDelete(body)
		if err != nil {
			return nil, err
		}
		return &Record{Op: op, Delete: del}, nil
	default:
		return nil, fmt.Errorf("unknown op type %d", op)
	}
}

func writeString(buf []byte, s string) []byte {
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(s)))
	buf = append(buf, s...)
	return buf
}

func readString(buf []byte) (string, []byte, error) {
	if len(buf) < 4 {
		return "", nil, errors.New("string: truncated length")
	}
	n := binary.BigEndian.Uint32(buf[:4])
	buf = buf[4:]
	if uint32(len(buf)) < n {
		return "", nil, errors.New("string: truncated body")
	}
	return string(buf[:n]), buf[n:], nil
}

func marshalCreateColl(c *CreateColl) ([]byte, error) {
	if c == nil {
		return nil, errors.New("CreateColl: nil")
	}
	buf := make([]byte, 0, 1+4+len(c.Name)+4+1+4+4)
	buf = append(buf, byte(OpCreateColl))
	buf = writeString(buf, c.Name)
	buf = binary.BigEndian.AppendUint32(buf, uint32(c.Dim))
	buf = append(buf, c.Metric)
	buf = binary.BigEndian.AppendUint32(buf, uint32(c.M))
	buf = binary.BigEndian.AppendUint32(buf, uint32(c.EfConstruction))
	return buf, nil
}

func unmarshalCreateColl(body []byte) (*CreateColl, error) {
	name, body, err := readString(body)
	if err != nil {
		return nil, fmt.Errorf("CreateColl name: %w", err)
	}
	if len(body) < 4+1+4+4 {
		return nil, errors.New("CreateColl: truncated fixed fields")
	}
	c := &CreateColl{Name: name}
	c.Dim = int32(binary.BigEndian.Uint32(body[:4]))
	c.Metric = body[4]
	c.M = int32(binary.BigEndian.Uint32(body[5:9]))
	c.EfConstruction = int32(binary.BigEndian.Uint32(body[9:13]))
	return c, nil
}

func marshalDropColl(d *DropColl) ([]byte, error) {
	if d == nil {
		return nil, errors.New("DropColl: nil")
	}
	buf := make([]byte, 0, 1+4+len(d.Name))
	buf = append(buf, byte(OpDropColl))
	buf = writeString(buf, d.Name)
	return buf, nil
}

func unmarshalDropColl(body []byte) (*DropColl, error) {
	name, _, err := readString(body)
	if err != nil {
		return nil, fmt.Errorf("DropColl name: %w", err)
	}
	return &DropColl{Name: name}, nil
}

func marshalInsert(ins *Insert) ([]byte, error) {
	if ins == nil {
		return nil, errors.New("Insert: nil")
	}
	dim := len(ins.Vector)
	buf := make([]byte, 0, 1+4+len(ins.Collection)+4+4+4*dim)
	buf = append(buf, byte(OpInsert))
	buf = writeString(buf, ins.Collection)
	buf = binary.BigEndian.AppendUint32(buf, ins.ID)
	buf = binary.BigEndian.AppendUint32(buf, uint32(dim))
	for _, v := range ins.Vector {
		buf = binary.BigEndian.AppendUint32(buf, math.Float32bits(v))
	}
	return buf, nil
}

func unmarshalInsert(body []byte) (*Insert, error) {
	coll, body, err := readString(body)
	if err != nil {
		return nil, fmt.Errorf("Insert collection: %w", err)
	}
	if len(body) < 8 {
		return nil, errors.New("Insert: truncated id/dim")
	}
	id := binary.BigEndian.Uint32(body[:4])
	dim := binary.BigEndian.Uint32(body[4:8])
	body = body[8:]
	if uint32(len(body)) < 4*dim {
		return nil, errors.New("Insert: truncated vector")
	}
	vec := make([]float32, dim)
	for i := uint32(0); i < dim; i++ {
		vec[i] = math.Float32frombits(binary.BigEndian.Uint32(body[4*i : 4*i+4]))
	}
	return &Insert{Collection: coll, ID: id, Vector: vec}, nil
}

func marshalDelete(d *Delete) ([]byte, error) {
	if d == nil {
		return nil, errors.New("Delete: nil")
	}
	buf := make([]byte, 0, 1+4+len(d.Collection)+4)
	buf = append(buf, byte(OpDelete))
	buf = writeString(buf, d.Collection)
	buf = binary.BigEndian.AppendUint32(buf, d.ID)
	return buf, nil
}

func unmarshalDelete(body []byte) (*Delete, error) {
	coll, body, err := readString(body)
	if err != nil {
		return nil, fmt.Errorf("Delete collection: %w", err)
	}
	if len(body) < 4 {
		return nil, errors.New("Delete: truncated id")
	}
	return &Delete{Collection: coll, ID: binary.BigEndian.Uint32(body[:4])}, nil
}
