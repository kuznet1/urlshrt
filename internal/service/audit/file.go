package audit

import (
	"encoding/json"
	"fmt"
	"github.com/kuznet1/urlshrt/internal/model"
	"os"
	"time"
)

type FileAudit struct {
	file *os.File
}

func NewFile(fname string) (*FileAudit, error) {
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create audit logs file %s: %w", fname, err)
	}
	return &FileAudit{
		file: f,
	}, nil
}

func (a *FileAudit) OnAuditEvt(userID int, action model.AuditAction, url string) error {
	return json.NewEncoder(a.file).Encode(model.AuditEvent{
		TS:     time.Now().Unix(),
		UserID: userID,
		Action: action,
		URL:    url,
	})
}

func (a *FileAudit) Close() error {
	return a.file.Close()
}
