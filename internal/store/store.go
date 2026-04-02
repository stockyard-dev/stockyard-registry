package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct { db *sql.DB }

type Package struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Type         string   `json:"type"`
	Size         int      `json:"size"`
	Digest       string   `json:"digest"`
	CreatedAt    string   `json:"created_at"`
}

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	dsn := filepath.Join(dataDir, "registry.db") + "?_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS packages (
			id TEXT PRIMARY KEY,\n\t\t\tname TEXT DEFAULT '',\n\t\t\tversion TEXT DEFAULT '',\n\t\t\ttype TEXT DEFAULT 'docker',\n\t\t\tsize INTEGER DEFAULT 0,\n\t\t\tdigest TEXT DEFAULT '',
			created_at TEXT DEFAULT (datetime('now'))
		)`)
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }

func genID() string { return fmt.Sprintf("%d", time.Now().UnixNano()) }

func (d *DB) Create(e *Package) error {
	e.ID = genID()
	e.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err := d.db.Exec(`INSERT INTO packages (id, name, version, type, size, digest, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.ID, e.Name, e.Version, e.Type, e.Size, e.Digest, e.CreatedAt)
	return err
}

func (d *DB) Get(id string) *Package {
	row := d.db.QueryRow(`SELECT id, name, version, type, size, digest, created_at FROM packages WHERE id=?`, id)
	var e Package
	if err := row.Scan(&e.ID, &e.Name, &e.Version, &e.Type, &e.Size, &e.Digest, &e.CreatedAt); err != nil {
		return nil
	}
	return &e
}

func (d *DB) List() []Package {
	rows, err := d.db.Query(`SELECT id, name, version, type, size, digest, created_at FROM packages ORDER BY created_at DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var result []Package
	for rows.Next() {
		var e Package
		if err := rows.Scan(&e.ID, &e.Name, &e.Version, &e.Type, &e.Size, &e.Digest, &e.CreatedAt); err != nil {
			continue
		}
		result = append(result, e)
	}
	return result
}

func (d *DB) Delete(id string) error {
	_, err := d.db.Exec(`DELETE FROM packages WHERE id=?`, id)
	return err
}

func (d *DB) Count() int {
	var n int
	d.db.QueryRow(`SELECT COUNT(*) FROM packages`).Scan(&n)
	return n
}
