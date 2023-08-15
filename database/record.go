package database

import "github.com/google/uuid"

type Status string

const (
	S_PENDING Status = "pending"
	S_READY   Status = "ready"
	S_FAILED  Status = "failed"
)

type Record struct {
	Id       uuid.UUID
	Filename string
	Status   Status
	MimeType string
	Path     string
	JSON     string
}

func New(id uuid.UUID) Record {
	return Record{
		Id:     id,
		Status: S_PENDING,
	}
}
