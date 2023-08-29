package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type InMemoryDb struct {
	receipts map[uuid.UUID]Receipt
}

func NewInMemoryDb() *InMemoryDb {
	return &InMemoryDb{
		receipts: make(map[uuid.UUID]Receipt),
	}
}

func (db *InMemoryDb) Get(id uuid.UUID) (Receipt, error) {
	if receipt, ok := db.receipts[id]; ok {
		return receipt, nil
	}
	return Receipt{}, errors.New("receipt not found")
}

func (db *InMemoryDb) Create(receipt Receipt) error {
	if _, ok := db.receipts[receipt.Id]; !ok {
		db.receipts[receipt.Id] = receipt
		return nil
	}
	return fmt.Errorf("receipt with uuid %v already exists", receipt.Id)
}

func (db *InMemoryDb) Update(receipt Receipt) error {
	if _, ok := db.receipts[receipt.Id]; ok {
		db.receipts[receipt.Id] = receipt
		return nil
	}
	return errors.New("receipt not found")
}
