package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/likeawizard/document-ai-demo/config"
)

var Instance DB

const (
	DRIVER_SQLITE    = "sqlite"
	DRIVER_IN_MEMORY = "inmemory"
)

type DB interface {
	Get(uuid.UUID) (Record, error)
	Create(Record) error
	Update(Record) error
}

func NewDataBase(cfg config.DbCfg) (DB, error) {
	switch cfg.Driver {
	case DRIVER_SQLITE:
		return NewSQLiteDb(cfg)
	case DRIVER_IN_MEMORY:
		return NewInMemoryDb(), nil
	default:
		return nil, fmt.Errorf("unsupported database driver %s", cfg.Driver)
	}
}