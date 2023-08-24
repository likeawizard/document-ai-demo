package database

import "github.com/google/uuid"

type Status string

const (
	S_PENDING Status = "pending"
	S_READY   Status = "ready"
	S_FAILED  Status = "failed"
)

type Record struct {
	Id       uuid.UUID `json:"id"`
	Filename string    `json:"filename"`
	Status   Status    `json:"status"`
	MimeType string    `json:"mime_type"`
	Path     string    `json:"path"`
	JSON     string    `json:"json_path"`
}

func New(id uuid.UUID) Record {
	return Record{
		Id:     id,
		Status: S_PENDING,
	}
}
