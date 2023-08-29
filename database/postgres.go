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

func (ps *PostgresDb) Get(id uuid.UUID) (Receipt, error) {
	r := Receipt{}
	sql := "SELECT id, filename, status, mime_type, path FROM receipts WHERE id = $1"
	row := ps.db.QueryRow(context.Background(), sql, id)
	err := row.Scan(&r.Id, &r.Filename, &r.Status, &r.MimeType, &r.Path)
	if err != nil {
		return r, fmt.Errorf("failed to retrieve receipt for id %s: %w", id, err)
	}
	return r, nil
}

func (ps *PostgresDb) Create(receipt Receipt) error {
	sql := "INSERT INTO receipts (id, filename, status, mime_type, path) VALUES($1, $2, $3, $4, $5)"
	_, err := ps.db.Exec(context.Background(), sql, receipt.Id, receipt.Filename, receipt.Status, receipt.MimeType, receipt.Path)
	if err != nil {
		return fmt.Errorf("failed to create new receipt with id %s: %w", receipt.Id, err)
	}
	return nil
}

func (ps *PostgresDb) Update(receipt Receipt) error {
	sql := "UPDATE receipts SET filename=$1, status=$2, mime_type=$3, path=$4 WHERE id=$5"
	_, err := ps.db.Exec(context.Background(), sql, receipt.Filename, receipt.Status, receipt.MimeType, receipt.Path, receipt.Id)
	if err != nil {
		return fmt.Errorf("failed to create new receipt with id %s: %w", receipt.Id, err)
	}
	return nil
}
