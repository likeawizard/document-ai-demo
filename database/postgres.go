package database

import (
	"context"
	"fmt"
	"strings"

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
	sql := `SELECT r.id, r.filename, r.status, r.mime_type, r.path, t.name FROM receipts r
		LEFT JOIN tags_to_receipts rel ON rel.receipt_id = r.id
		LEFT JOIN tags t ON t.id = rel.tag_id
		WHERE r.id = $1`
	rows, err := ps.db.Query(context.Background(), sql, id)
	if err != nil {
		return r, fmt.Errorf("failed to retrieve receipt for id %s: %w", id, err)
	}
	for rows.Next() {
		tag := ""
		err := rows.Scan(&r.Id, &r.Filename, &r.Status, &r.MimeType, &r.Path, &tag)
		if err != nil {
			continue
		}
		r.Tags = append(r.Tags, tag)
	}

	if err != nil {
		return r, fmt.Errorf("failed to retrieve receipt for id %s: %w", id, err)
	}
	return r, nil
}

// TODO/FIX: will return correct receipts but will omit receipt tags that were not passed as parameters
// might need to do two joins of tags- one for selecting receipts one for attaching all tags to selected receipts
func (ps *PostgresDb) GetByTags(tags []string) ([]Receipt, error) {
	receipts := make([]Receipt, 0)
	placeholders := make([]string, 0)
	args := make([]interface{}, 0)
	for i, v := range tags {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, v)
	}

	sql := fmt.Sprintf(`SELECT r.id, r.filename, r.status, r.mime_type, r.path, t.name FROM receipts r
		LEFT JOIN tags_to_receipts rel ON rel.receipt_id = r.id
		LEFT JOIN tags t ON t.id = rel.tag_id
		WHERE t.name IN (%s)`, strings.Join(placeholders, ", "))

	rows, err := ps.db.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve receipts with tags %s: %w", tags, err)
	}

	count := -1
	for rows.Next() {
		var tmpReceipt Receipt
		tag := ""
		err := rows.Scan(&tmpReceipt.Id, &tmpReceipt.Filename, &tmpReceipt.Status, &tmpReceipt.MimeType, &tmpReceipt.Path, &tag)
		if err != nil {
			continue
		}

		if count < 0 || tmpReceipt.Id != receipts[count].Id {
			receipts = append(receipts, tmpReceipt)
			count++
		}
		receipts[count].Tags = append(receipts[count].Tags, tag)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve receipts with tags %s: %w", tags, err)
	}

	return receipts, nil
}

func (ps *PostgresDb) Create(receipt Receipt) error {
	sql := "INSERT INTO receipts (id, filename, status, mime_type, path) VALUES($1, $2, $3, $4, $5)"
	_, err := ps.db.Exec(context.Background(), sql, receipt.Id, receipt.Filename, receipt.Status, receipt.MimeType, receipt.Path)
	if err != nil {
		return fmt.Errorf("failed to create new receipt with id %s: %w", receipt.Id, err)
	}
	ps.setTags(receipt.Id, receipt.Tags)
	return nil
}

func (ps *PostgresDb) Update(receipt Receipt) error {
	sql := "UPDATE receipts SET filename=$1, status=$2, mime_type=$3, path=$4 WHERE id=$5"
	_, err := ps.db.Exec(context.Background(), sql, receipt.Filename, receipt.Status, receipt.MimeType, receipt.Path, receipt.Id)
	if err != nil {
		return fmt.Errorf("failed to create new receipt with id %s: %w", receipt.Id, err)
	}
	ps.setTags(receipt.Id, receipt.Tags)
	return nil
}

// TODO: no error checking. reduce number of queries. maybe use a stored procedure for this
func (ps *PostgresDb) setTags(id uuid.UUID, tags []string) error {
	if len(tags) < 1 {
		return nil
	}
	deleteSql := `DELETE FROM tags_to_receipts WHERE receipt_id = $1`
	ps.db.Exec(context.Background(), deleteSql, id)

	inserts := make([]string, 0)
	rels := make([]string, 0)
	args := make([]interface{}, 0)
	for i, v := range tags {
		inserts = append(inserts, fmt.Sprintf("($%d)", i+1))
		rels = append(rels, fmt.Sprintf("$%d", i+2))
		args = append(args, v)
	}
	insertSql := fmt.Sprintf(`INSERT INTO tags (name) VALUES %s on CONFLICT DO NOTHING`, strings.Join(inserts, ", "))
	ps.db.Exec(context.Background(), insertSql, args...)

	linkSql := fmt.Sprintf(`INSERT INTO tags_to_receipts (tag_id, receipt_id)
	SELECT t.id, $1
	FROM tags t
	WHERE t.name IN (%s);`, strings.Join(rels, ", "))
	args = append([]interface{}{id}, args...)
	ps.db.Exec(context.Background(), linkSql, args...)
	return nil
}
