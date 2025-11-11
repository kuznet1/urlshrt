package model

type AuditAction string

const ActionShorten AuditAction = "shorten"
const ActionFollow AuditAction = "follow"

type AuditEvent struct {
	TS     int64       `json:"ts"`
	Action AuditAction `json:"action"`
	UserID int         `json:"user_id"`
	URL    string      `json:"url"`
}
