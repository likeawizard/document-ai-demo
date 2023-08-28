package database

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/likeawizard/document-ai-demo/config"
)

type PostgresDb struct {
	db *pgxpool.Pool
}

func NewPostgres(cfg config.DbCfg) (*PostgresDb, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &PostgresDb{db: conn}, nil
}

func (ps *PostgresDb) Get(id uuid.UUID) (Record, error) {
	r := Record{}
	sql := "SELECT uuid, filename, status, mime_type, path, json FROM records WHERE uuid = $1"
	row := ps.db.QueryRow(context.Background(), sql, id)
	err := row.Scan(&r.Id, &r.Filename, &r.Status, &r.MimeType, &r.Path, &r.JSON)
	if err != nil {
		return r, fmt.Errorf("failed to retrieve record for uuid %s: %w", id, err)
	}
	return r, nil
}

func (ps *PostgresDb) Create(record Record) error {
	sql := "INSERT INTO records (uuid, filename, status, mime_type, path, json) VALUES($1, $2, $3, $4, $5, $6)"
	_, err := ps.db.Exec(context.Background(), sql, record.Id, record.Filename, record.Status, record.MimeType, record.Path, record.JSON)
	if err != nil {
		return fmt.Errorf("failed to create new record with uuid %s: %w", record.Id, err)
	}
	return nil
}

func (ps *PostgresDb) Update(record Record) error {
	sql := "UPDATE records SET filename=$1, status=$2, mime_type=$3, path=$4, json=$5 WHERE uuid=$6"
	_, err := ps.db.Exec(context.Background(), sql, record.Filename, record.Status, record.MimeType, record.Path, record.JSON, record.Id)
	if err != nil {
		return fmt.Errorf("failed to create new record with uuid %s: %w", record.Id, err)
	}
	return nil
}
