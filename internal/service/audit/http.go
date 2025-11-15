package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kuznet1/urlshrt/internal/model"
	"net/http"
	"time"
)

// URLAudit forwards audit events to a remote HTTP endpoint.
type URLAudit struct {
	url string
}

// NewURLAudit creates an HTTP audit subscriber that POSTs JSON events to the configured URL.
func NewURLAudit(url string) *URLAudit {
	return &URLAudit{
		url: url,
	}
}

// OnAuditEvt sends the given event to the configured HTTP endpoint.
// It implements the AuditSubscriber interface.
func (a *URLAudit) OnAuditEvt(ctx context.Context, userID int, action model.AuditAction, url string) error {
	data, err := json.Marshal(model.AuditEvent{
		TS:     time.Now().Unix(),
		UserID: userID,
		Action: action,
		URL:    url,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}
