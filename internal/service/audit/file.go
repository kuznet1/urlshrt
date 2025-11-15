package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kuznet1/urlshrt/internal/model"
	"os"
	"time"
)

// FileAudit writes audit events to an io.Writer (typically a file).
type FileAudit struct {
	file *os.File
}

// NewFile performs a public package operation. Top-level handler/function.
func NewFile(fname string) (*FileAudit, error) {
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create audit logs file %s: %w", fname, err)
	}
	return &FileAudit{
		file: f,
	}, nil
}

// OnAuditEvt writes the event to the underlying writer in JSON format.
// It implements the AuditSubscriber interface.
func (a *FileAudit) OnAuditEvt(ctx context.Context, userID int, action model.AuditAction, url string) error {
	return json.NewEncoder(a.file).Encode(model.AuditEvent{
		TS:     time.Now().Unix(),
		UserID: userID,
		Action: action,
		URL:    url,
	})
}

// Close is a method that provides public behavior for the corresponding type.
func (a *FileAudit) Close() error {
	return a.file.Close()
}
