package database

// Can cause problems with cgo during compile time
import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/likeawizard/document-ai-demo/config"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	db *sql.DB
}

func (sqlite *SQLite) init() error {
	var createRecordsTable = `
	CREATE TABLE IF NOT EXISTS records (
		uuid TEXT PRIMARY KEY,
		filename TEXT,
		status TEXT,
		mime_type TEXT,
		path TEXT,
		json TEXT
	);
	`
	_, err := sqlite.db.Exec(createRecordsTable)
	if err != nil {
		return fmt.Errorf("failed to initialize sqlite: %w", err)
	}
	return nil
}

func NewSQLiteDb(cfg config.DbCfg) (*SQLite, error) {
	db, err := sql.Open("sqlite3", cfg.Name)
	if err != nil {
		return nil, err
	}
	sqlite := SQLite{
		db: db,
	}

	err = sqlite.init()
	if err != nil {
		return nil, err
	}

	return &sqlite, nil
}

func (sqlite *SQLite) Get(id uuid.UUID) (Record, error) {
	r := Record{}
	sql := "SELECT uuid, filename, status, mime_type, path, json FROM records WHERE uuid = ?"
	row := sqlite.db.QueryRow(sql, id)
	err := row.Scan(&r.Id, &r.Filename, &r.Status, &r.MimeType, &r.Path, &r.JSON)
	if err != nil {
		return r, fmt.Errorf("failed to retrieve record for uuid %s: %w", id, err)
	}
	return r, nil
}

func (sqlite *SQLite) Create(record Record) error {
	sql := "INSERT INTO records (uuid, filename, status, mime_type, path, json) VALUES(?, ?, ?, ?, ?, ?)"
	_, err := sqlite.db.Exec(sql, record.Id, record.Filename, record.Status, record.MimeType, record.Path, record.JSON)
	if err != nil {
		return fmt.Errorf("failed to create new record with uuid %s: %w", record.Id, err)
	}
	return nil
}

func (sqlite *SQLite) Update(record Record) error {
	sql := "UPDATE records SET filename=?, status=?, mime_type=?, path=?, json=? WHERE uuid=?"
	_, err := sqlite.db.Exec(sql, record.Filename, record.Status, record.MimeType, record.Path, record.JSON, record.Id)
	if err != nil {
		return fmt.Errorf("failed to create new record with uuid %s: %w", record.Id, err)
	}
	return nil
}
