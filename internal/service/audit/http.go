package audit

import (
	"bytes"
	"encoding/json"
	"github.com/kuznet1/urlshrt/internal/model"
	"net/http"
	"time"
)

type URLAudit struct {
	url string
}

func NewURLAudit(url string) *URLAudit {
	return &URLAudit{
		url: url,
	}
}

func (a *URLAudit) OnAuditEvt(userID int, action model.AuditAction, url string) error {
	data, err := json.Marshal(model.AuditEvent{
		TS:     time.Now().Unix(),
		UserID: userID,
		Action: action,
		URL:    url,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(data))
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}
