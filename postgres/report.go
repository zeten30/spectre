package postgres

import (
	"context"
	"time"

	"github.com/zeten30/spectre"
)

type dbReport struct {
	PasteID string
	Count   int

	conn *conn
	ctx  context.Context
}

func (r *dbReport) GetPasteID() spectre.PasteID {
	return spectre.PasteID(r.PasteID)
}

func (r *dbReport) GetCount() int {
	return r.Count
}

func (r *dbReport) GetCreationTime() time.Time {
	return time.Now()
}

func (r *dbReport) GetModificationTime() time.Time {
	return time.Now()
}

func (r *dbReport) Commit() error {
	// TODO(DH)
	// can these even be committed?
	return nil
}

func (r *dbReport) Erase() error {
	_, err := r.conn.db.ExecContext(r.ctx, `DELETE FROM paste_reports WHERE paste_id = $1`, r.PasteID)
	return err
}
