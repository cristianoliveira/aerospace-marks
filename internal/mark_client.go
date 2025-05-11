package internal

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Mark struct {
	ID        int
	Mark      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MarkClient struct {
	db *sql.DB
}

func NewMarkClient(dbPath string) (*MarkClient, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	client := &MarkClient{db: db}
	if err := client.createTableIfNotExists(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *MarkClient) createTableIfNotExists() error {
	query := `
	CREATE TABLE IF NOT EXISTS marks (
		id INTEGER PRIMARY KEY,
		mark TEXT,
		created_at DATETIME,
		updated_at DATETIME
	);
	`
	_, err := c.db.Exec(query)
	return err
}

func (c *MarkClient) AddMark(id int, mark string) error {
	query := `
	INSERT INTO marks (id, mark, created_at, updated_at)
	VALUES (?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		mark=excluded.mark,
		updated_at=excluded.updated_at;
	`
	_, err := c.db.Exec(query, id, mark, time.Now(), time.Now())
	return err
}

func (c *MarkClient) Close() error {
	return c.db.Close()
}
