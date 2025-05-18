package storage

import (
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type MarkStorage interface {
	// AddMark adds a mark to the database
	AddMark(id string, mark string) error
	// GetMarks returns all marks in the database
	GetMarks() ([]Mark, error)
	// GetMarksByWindowID returns all marks for a given window ID
	GetMarksByWindowID(id string) ([]Mark, error)
	// GetWindowIDByMark returns the window ID for a given mark
	GetWindowIDByMark(mark string) (string, error)
	// ReplaceAllMarks replaces all marks for a window with a new mark
	ReplaceAllMarks(id string, mark string) (bool, error)
	// Close closes the database connection
	Close() error
}

type MarkStorageClient struct {
	storage StorageDbClient
}

func NewMarkClient() (*MarkStorageClient, error) {
	storageClient, err := DefaultConnector.Connect()
	if err != nil {
		return nil, err
	}

	client := &MarkStorageClient{
		storage: storageClient,
	}

	return client, nil
}

func (c *MarkStorageClient) AddMark(id string, mark string) error {
	query := strings.TrimSpace(`
	INSERT INTO marks (window_id, mark) VALUES (?, ?)
	`)
	_, err := c.storage.Execute(query, id, mark)
	return err
}

func (c *MarkStorageClient) GetMarks() ([]Mark, error) {
	query := `
	SELECT window_id, mark
	FROM marks
	`
	marks, err := c.storage.QueryAll(query)
	return marks, err
}

func (c *MarkStorageClient) GetMarksByWindowID(id string) ([]Mark, error) {
	query := `
	SELECT window_id, mark
	FROM marks
	WHERE window_id = ?
	`
	marks, err := c.storage.QueryAll(query, id)
	if err != nil {
		return nil, err
	}
	return marks, nil
}

// Get window ID by mark
//
// This function will return the first window ID that matches the mark
// If multiple window IDs match the mark, it will return the first one found
func (c *MarkStorageClient) GetWindowIDByMark(markI string) (string, error) {
	query := "SELECT * FROM marks WHERE mark = ?"

	markedWindow, err := c.storage.QueryOne(query, markI)
	if err != nil {
		return "", err
	}

	return markedWindow.WindowID, nil
}

// ReplaceAllMarks replaces all marks for a window with a new mark
// This function will delete all marks for the specified window ID and
// then add the new mark
func (c *MarkStorageClient) ReplaceAllMarks(id string, mark string) (int64, error) {
	// Delete all marks for the window
	query := strings.TrimSpace(`
	DELETE FROM marks WHERE mark = ?
	`)

	res, err := c.storage.Execute(query, mark)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	// Add the new mark
	return rowsAffected, c.AddMark(id, mark)
}

func (c *MarkStorageClient) Close() error {
	return c.storage.Close()
}
