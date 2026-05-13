package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/zhangchuqi1998/velosearch/internal/collection"
	"github.com/zhangchuqi1998/velosearch/internal/storage"
)

// ApplyWALRecord applies a single WAL record to the Manager.
// Used by both startup replay and runtime apply.
//
// Replay is idempotent: re-applying a record that already took effect
// (e.g. Create on existing name, Delete on missing id) logs a warning
// and returns nil so replay can continue past inconsistencies caused
// by client-side retries.
func ApplyWALRecord(mgr *collection.Manager, rec *storage.Record) error {
	switch rec.Op {
	case storage.OpCreateColl:
		c := rec.CreateColl
		err := mgr.Create(collection.Config{
			Name:           c.Name,
			Dim:            int(c.Dim),
			Metric:         collection.Metric(c.Metric),
			M:              int(c.M),
			EfConstruction: int(c.EfConstruction),
		})
		if errors.Is(err, collection.ErrAlreadyExists) {
			slog.Warn("wal replay: collection already exists, skipping", "name", c.Name)
			return nil
		}
		return err
	case storage.OpDropColl:
		err := mgr.Drop(rec.DropColl.Name)
		if errors.Is(err, collection.ErrNotFound) {
			slog.Warn("wal replay: drop on missing collection, skipping", "name", rec.DropColl.Name)
			return nil
		}
		return err
	case storage.OpInsert:
		col, err := mgr.Get(rec.Insert.Collection)
		if err != nil {
			return err
		}
		col.Index.Insert(rec.Insert.ID, rec.Insert.Vector)
		return nil
	case storage.OpDelete:
		col, err := mgr.Get(rec.Delete.Collection)
		if err != nil {
			return err
		}
		if err := col.Index.Delete(rec.Delete.ID); err != nil {
			slog.Warn("wal replay: delete failed, skipping",
				"collection", rec.Delete.Collection, "id", rec.Delete.ID, "err", err)
		}
		return nil
	}
	return fmt.Errorf("unknown op type %d", rec.Op)
}
