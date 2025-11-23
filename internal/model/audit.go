package model

type AuditAction string

// ActionShorten is a public package constant used for configuration or external access.
const ActionShorten AuditAction = "shorten"

// ActionFollow is a public package constant used for configuration or external access.
const ActionFollow AuditAction = "follow"

// AuditEvent is a public struct of the package. It exposes the core data for this project.
type AuditEvent struct {
	TS     int64       `json:"ts"`
	Action AuditAction `json:"action"`
	UserID int         `json:"user_id"`
	URL    string      `json:"url"`
}
