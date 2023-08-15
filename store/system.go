package store

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/likeawizard/document-ai-demo/config"
)

type SystemStore struct {
	base string
}

func (ss *SystemStore) GetPath(filename string) (string, error) {
	if filename == "" {
		return "", errors.New("empty filename")
	}
	return filepath.Join(ss.base, filename), nil
}

func NewSystemStore(cfg config.StorageCfg) *SystemStore {
	return &SystemStore{base: cfg.Location}
}

func (ss *SystemStore) Get(filename string) (io.ReadCloser, error) {
	path, err := ss.GetPath(filename)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func (ss *SystemStore) Store(filename string, r io.Reader) error {
	path, err := ss.GetPath(filename)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	br := bufio.NewReader(r)
	_, err = br.WriteTo(f)
	return err
}
