package database

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/likeawizard/document-ai-demo/config"
)

const (
	DRIVER_POSTGRES  = "postgres"
	DRIVER_IN_MEMORY = "inmemory"
)

type DB interface {
	Get(uuid.UUID) (Receipt, error)
	GetByTags([]string) ([]Receipt, error)
	Create(Receipt) error
	Update(Receipt) error
}

func NewDataBase(cfg config.DbCfg) (DB, error) {
	switch cfg.Driver {
	case DRIVER_POSTGRES:
		return NewPostgres(cfg)
	case DRIVER_IN_MEMORY:
		return NewInMemoryDb(), nil
	default:
		return nil, fmt.Errorf("unsupported database driver %s", cfg.Driver)
	}
}
