package store

import (
	"fmt"
	"io"

	"github.com/likeawizard/document-ai-demo/config"
)

const (
	DRIVER_FS     = "os"
	DRIVER_GCLOUD = "gcloud"
)

type FileStore interface {
	Get(string) (io.ReadCloser, error)
	Store(string, io.Reader) error
	GetURL(string) (string, error)
}

var File FileStore

func NewFileStore(cfg config.StorageCfg) (FileStore, error) {
	switch cfg.Driver {
	case DRIVER_FS:
		return NewSystemStore(cfg), nil
	case DRIVER_GCLOUD:
		return NewGCloudStore(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported file store driver %s", cfg.Driver)
	}
}
