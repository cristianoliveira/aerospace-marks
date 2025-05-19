package storage

import (
	"testing"

	"github.com/cristianoliveira/aerospace-marks/internal/logger"
)

// TestMarksStorageClient tests the MarksStorageClient functionality.
type MockConnector struct {
	Client StorageDbClient
}
func (t *MockConnector) Connect() (StorageDbClient, error) {
	return t.Client, nil
}

// Map string, any 
type MockMarkStorage struct {
	RecordArgs []string
	Marks		   []Mark
	Mark       *Mark
	DbResult   DbResult
}

func (m *MockMarkStorage) QueryAll(query string, args ...any) ([]Mark, error) {
	m.RecordArgs = append(m.RecordArgs, query)
	return m.Marks, nil
}
func (m *MockMarkStorage) QueryOne(query string, args ...any) (*Mark, error) {
	m.RecordArgs = append(m.RecordArgs, query)
	return m.Mark, nil
}
func (m *MockMarkStorage) Execute(query string, args ...any) (DbResult, error) {
	m.RecordArgs = append(m.RecordArgs, query)
	return m.DbResult, nil 
}
func (m *MockMarkStorage) Close() error {
	m.RecordArgs = append(m.RecordArgs, "close")
	return nil
}

func TestMarksStorageClient(t *testing.T) {
	logger.SetDefaultLogger(&logger.EmptyLogger{})

	t.Run("Test NewMarkClient", func(t *testing.T) {
		DefaultConnector = &MockConnector{
			Client: &MockMarkStorage{},
		}

		client, err := NewMarkClient()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if client == nil {
			t.Fatal("expected client to be non-nil")
		}
	})

	t.Run("Test AddMark", func(t *testing.T) {
		DefaultConnector = &MockConnector{
			Client: &MockMarkStorage{},
		}

		client, err := NewMarkClient()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer client.Close()

		err = client.AddMark("window1", "mark1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(client.(*MarkStorageClient).storage.(*MockMarkStorage).RecordArgs) != 1 {
			t.Fatalf("expected 1 query, got %d", len(client.(*MarkStorageClient).storage.(*MockMarkStorage).RecordArgs))
		}
		expectedQuery := "INSERT INTO marks (window_id, mark) VALUES (?, ?)"
		if client.(*MarkStorageClient).storage.(*MockMarkStorage).RecordArgs[0] != expectedQuery {
			t.Fatalf("expected query %s, got %s", expectedQuery, client.(*MarkStorageClient).storage.(*MockMarkStorage).RecordArgs[0])
		}
	})

	t.Run("Test GetMarks", func(t *testing.T) {
		DefaultConnector = &MockConnector{
			Client: &MockMarkStorage{
				Marks: []Mark{
					{WindowID: "window1", Mark: "mark1"},
				},
			},
		}

		client, err := NewMarkClient()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer client.Close()

		marks, err := client.GetMarks()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if marks == nil {
			t.Fatal("expected marks to be non-nil")
		}

		if len(client.(*MarkStorageClient).storage.(*MockMarkStorage).RecordArgs) != 1 {
			t.Fatalf("expected 1 query, got %d", len(client.(*MarkStorageClient).storage.(*MockMarkStorage).RecordArgs))
		}
		expectedQuery := "SELECT window_id, mark FROM marks"
		if client.(*MarkStorageClient).storage.(*MockMarkStorage).RecordArgs[0] != expectedQuery {
			t.Fatalf("expected query %s, got %s", expectedQuery, client.(*MarkStorageClient).storage.(*MockMarkStorage).RecordArgs[0])
		}
	})
}
