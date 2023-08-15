package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type InMemoryDb struct {
	records map[uuid.UUID]Record
}

func NewInMemoryDb() *InMemoryDb {
	return &InMemoryDb{
		records: make(map[uuid.UUID]Record),
	}
}

func (db *InMemoryDb) Get(id uuid.UUID) (Record, error) {
	if record, ok := db.records[id]; ok {
		return record, nil
	}
	return Record{}, errors.New("record not found")
}

func (db *InMemoryDb) Create(record Record) error {
	if _, ok := db.records[record.Id]; !ok {
		db.records[record.Id] = record
		return nil
	}
	return fmt.Errorf("record with uuid %v already exists", record.Id)
}

func (db *InMemoryDb) Update(record Record) error {
	if _, ok := db.records[record.Id]; ok {
		db.records[record.Id] = record
		return nil
	}
	return errors.New("record not found")
}
