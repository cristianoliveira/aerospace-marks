package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Mark struct {
	windowID  string
	Mark      string
}

type MarkClient struct {
	db *sql.DB
}

func NewMarkClient() (*MarkClient, error) {
	// Create the directory if it doesn't exist
	dir := fmt.Sprintf("%s/.local/state/aerospace-marks", os.Getenv("HOME"))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	dbPath := fmt.Sprintf("%s/.local/state/aerospace-marks/storage.db", os.Getenv("HOME"))
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
	// Marks are unique to a window, so we can use the window_id as a unique key
	// but one window may have multiple marks
	query := `
	CREATE TABLE IF NOT EXISTS marks (
		window_id TEXT,
		mark TEXT,
		constraint marks_pk PRIMARY KEY (window_id, mark)
	);
	`
	_, err := c.db.Exec(query)
	return err
}

func (c *MarkClient) AddMark(id string, mark string) error {
	query := `
	INSERT INTO marks (window_id, mark)
	VALUES (?, ?)
	`
	_, err := c.db.Exec(query, id, mark, time.Now(), time.Now())
	return err
}

func (c *MarkClient) GetMarks() ([]Mark, error) {
	query := `
	SELECT window_id, mark
	FROM marks
	`
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var marks []Mark
	for rows.Next() {
		var mark Mark
		if err := rows.Scan(&mark.windowID, &mark.Mark); err != nil {
			return nil, err
		}
		marks = append(marks, mark)
	}

	return marks, nil
}

func (c *MarkClient) GetMarksByWindowID(id string) ([]Mark, error) {
	query := `
	SELECT window_id, mark
	FROM marks
	WHERE window_id = ?
	`
	rows, err := c.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var marks []Mark
	for rows.Next() {
		var mark Mark
		if err := rows.Scan(&mark.windowID, &mark.Mark); err != nil {
			return nil, err
		}
		marks = append(marks, mark)
	}

	return marks, nil
}

// Get window ID by mark
//
// This function will return the first window ID that matches the mark
// If multiple window IDs match the mark, it will return the first one found
func (c *MarkClient) GetWindowIDByMark(mark string) (string, error) {
	query := `
	SELECT window_id
	FROM marks
	WHERE mark = ?
	`
	row := c.db.QueryRow(query, mark)

	var windowID string
	if err := row.Scan(&windowID); err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No rows found
		}
		return "", err
	}

	return windowID, nil
}

// ReplaceAllMarks replaces all marks for a window with a new mark
// This function will delete all marks for the specified window ID and
// then add the new mark
// It returns true if marks were deleted, false if no marks were found
func (c *MarkClient) ReplaceAllMarks(id string, mark string) (bool, error) {
	// Delete all marks for the window
	query := `
	DELETE FROM marks
	WHERE window_id = ?
	`

	var hasDeleted bool
	res, err := c.db.Exec(query, id)
	if err != nil {
		return false, err
	}

	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		hasDeleted = false
	} else {
		hasDeleted = true
	}

	// Add the new mark
	return hasDeleted, c.AddMark(id, mark)
}

func (c *MarkClient) Close() error {
	return c.db.Close()
}
