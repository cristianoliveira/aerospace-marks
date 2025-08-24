package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cristianoliveira/aerospace-marks/internal/storage/db/queries"
	_ "github.com/mattn/go-sqlite3"
)

type MarkStorage interface {
	// AddMark adds a mark to the database
	AddMark(id string, mark string) error
	// GetMarks returns all marks in the database
	GetMarks() ([]queries.Mark, error)
	// GetMarksByWindowID returns all marks for a given window ID
	GetMarksByWindowID(id string) ([]queries.Mark, error)
	// GetWindowByMark returns the window for a given mark
	GetWindowByMark(mark string) (*queries.Mark, error)
	// GetWindowIDByMark returns the window ID for a given mark
	GetWindowIDByMark(mark string) (string, error)
	// ReplaceAllMarks replaces all marks for a window with a new mark
	ReplaceAllMarks(id string, mark string) (int64, error)
	// ToggleMark toggles a mark for a window
	ToggleMark(id string, mark string) error
	// DeleteByMark removes a mark from the database
	DeleteByMark(mark string) (int64, error)
	// DeleteByMark removes a mark from the database
	DeleteByWindow(windowId int) (int64, error)
	// DeleteAllMarks removes all marks from the database
	DeleteAllMarks() (int64, error)
	// Close closes the database connection
	Close() error
	// Client returns the storage client
	Client() StorageDbClient
}

type MarkStorageClient struct {
	storage StorageDbClient
	queries *queries.Queries
}

func NewMarkClient(storageClient StorageDbClient) (*MarkStorageClient, error) {
	// Initialize SQLC queries with the underlying database connection
	db := storageClient.GetDB()
	queriesClient := queries.New(db)

	client := &MarkStorageClient{
		storage: storageClient,
		queries: queriesClient,
	}

	return client, nil
}

func (c *MarkStorageClient) AddMark(id string, mark string) error {
	ctx := context.Background()
	return c.queries.AddMark(ctx, id, mark)
}

func (c *MarkStorageClient) GetMarks() ([]queries.Mark, error) {
	ctx := context.Background()
	return c.queries.GetAllMarks(ctx)
}

func (c *MarkStorageClient) GetMarksByWindowID(id string) ([]queries.Mark, error) {
	ctx := context.Background()
	return c.queries.GetMarksByWindowID(ctx, id)
}

// GetWindowByMark returns the window for a given mark
//
// This function will return the first window that matches the mark
// If multiple windows match the mark, it will error
func (c *MarkStorageClient) GetWindowByMark(mark string) (*queries.Mark, error) {
	ctx := context.Background()
	markedWindow, err := c.queries.GetWindowByMark(ctx, mark)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no window found for mark %s", mark)
		}
		return nil, err
	}

	return &markedWindow, nil
}

// Get window ID by mark
//
// This function will return the first window ID that matches the mark
// If multiple window IDs match the mark, it will return the first one found
func (c *MarkStorageClient) GetWindowIDByMark(markI string) (string, error) {
	ctx := context.Background()
	markedWindow, err := c.queries.GetWindowByMark(ctx, markI)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("no window found for mark %s", markI)
		}
		return "", err
	}

	return markedWindow.WindowID, nil
}

// ReplaceAllMarks replaces all marks for a window with a new mark
// This function will delete all marks for the specified window ID and
// then add the new mark
func (c *MarkStorageClient) ReplaceAllMarks(id string, mark string) (int64, error) {
	ctx := context.Background()

	// Delete all marks for the window
	res, err := c.queries.DeleteMarksByWindowIDOrMark(ctx, id, mark)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	err = c.AddMark(id, mark)
	if err != nil {
		return rowsAffected, err
	}

	return rowsAffected, nil
}

func (c *MarkStorageClient) Close() error {
	return c.storage.Close()
}

// Client returns the storage client
func (c *MarkStorageClient) Client() StorageDbClient {
	return c.storage
}

// ToggleMark toggles a mark for a window
// If the mark exists, it will be deleted
// If the mark does not exist, it will be added
func (c *MarkStorageClient) ToggleMark(id string, mark string) error {
	rowsAffected, err := c.DeleteByMark(mark)
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		// Mark was deleted
		return nil
	}

	// Mark was not deleted, so add it
	err = c.AddMark(id, mark)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAllMarks removes all marks from the database
func (c *MarkStorageClient) DeleteAllMarks() (int64, error) {
	ctx := context.Background()
	res, err := c.queries.DeleteAllMarks(ctx)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// DeleteByMark deletes a mark from the database
// This function will delete the mark from the database
func (c *MarkStorageClient) DeleteByMark(mark string) (int64, error) {
	ctx := context.Background()
	res, err := c.queries.DeleteByMark(ctx, mark)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

// DeleteByWindow deletes a mark from the database
// This function will delete the mark from the database
func (c *MarkStorageClient) DeleteByWindow(windowId int) (int64, error) {
	ctx := context.Background()
	// Convert int to string since our SQLC queries expect string
	windowIdStr := fmt.Sprintf("%d", windowId)
	res, err := c.queries.DeleteByWindow(ctx, windowIdStr)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	return rowsAffected, err
}
