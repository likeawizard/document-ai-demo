package store

import (
	"fmt"
	"io"

	"github.com/likeawizard/document-ai-demo/config"
)

const (
	DRIVER_FS = "os"
)

type FileStore interface {
	Get(string) (io.ReadCloser, error)
	Store(string, io.Reader) error
}

var File FileStore

func NewFileStore(cfg config.StorageCfg) (FileStore, error) {
	switch cfg.Driver {
	case DRIVER_FS:
		return NewSystemStore(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported file store driver %s", cfg.Driver)
	}
}
