package database

import (
	"fmt"

	"github.com/google/uuid"
)

type Status string

const (
	S_PENDING Status = "pending"
	S_READY   Status = "ready"
	S_FAILED  Status = "failed"
)

type Receipt struct {
	Id       uuid.UUID `json:"id"`
	Filename string    `json:"filename"`
	Status   Status    `json:"status"`
	MimeType string    `json:"mime_type"`
	Path     string    `json:"path"`
}

func New(id uuid.UUID) Receipt {
	return Receipt{
		Id:     id,
		Status: S_PENDING,
	}
}

func (r Receipt) GetJsonPath() string {
	return fmt.Sprintf("%s.json", r.Id)
}

func (r Receipt) GetExpensePath() string {
	return fmt.Sprintf("%s-expense.json", r.Id)
}