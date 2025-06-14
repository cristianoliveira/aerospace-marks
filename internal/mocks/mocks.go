package mocks

import (
	"encoding/json"
	"fmt"
	"os"

	aerospaceipc "github.com/cristianoliveira/aerospace-ipc"
	aerospacecli "github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"github.com/cristianoliveira/aerospace-marks/internal/aerospace"
	aerospacecli_mock "github.com/cristianoliveira/aerospace-marks/internal/mocks/aerospacecli"
	storage_mock "github.com/cristianoliveira/aerospace-marks/internal/mocks/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"go.uber.org/mock/gomock"
)

// This module contains a set of mock helpers for mocking AeroSpace socket connections
// and storage, easily used in unit tests.

func MockStorageDbClient(ctrl *gomock.Controller) (*storage_mock.MockStorageDbClient, storage.MarkStorage) {
	storageDbClient := storage_mock.NewMockStorageDbClient(ctrl)

	newStorage, err := storage.NewMarkClient(storageDbClient)
	if err != nil {
		ctrl.T.Errorf("failed to create new storage client: %v", err)
	}

	return storageDbClient, newStorage
}

func MockAerospaceConnection(ctrl *gomock.Controller) (
	*aerospacecli_mock.MockAeroSpaceConnection,
	aerospace.AerosSpaceMarkWindows,
) {
	mockAeroSpaceConnection := aerospacecli_mock.NewMockAeroSpaceConnection(ctrl)
	mockAeroSpaceConnetor := aerospacecli_mock.NewMockAeroSpaceConnector(ctrl)
	mockAeroSpaceConnetor.EXPECT().Connect().Return(mockAeroSpaceConnection, nil).Times(1)
	// mockAeroSpaceConnection.EXPECT().CloseConnection().Return(nil).Times(1)

	aerospacecli.SetDefaultConnector(mockAeroSpaceConnetor)

	aerospaceClient, err := aerospace.NewAeroSpaceClient()
	if err != nil {
		panic(fmt.Errorf("failed to create aerospace client: %w", err))
	}

	return mockAeroSpaceConnection, aerospaceClient
}

func MockStorageDbResult(ctrl *gomock.Controller, lastInsertId *int64, rowsAffected *int64) *storage_mock.MockDbResult {
	dbResult := storage_mock.NewMockDbResult(ctrl)
	if lastInsertId != nil {
		dbResult.EXPECT().
			LastInsertId().
			Return(*lastInsertId, nil).
			Times(1)
	}

	if rowsAffected != nil {
		dbResult.EXPECT().
			RowsAffected().
			Return(*rowsAffected, nil).
			Times(1)
	}

	return dbResult
}

func LoadMarksFixture(jsonFilePath string) ([]storage.Mark, error) {
	file, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	var marks []storage.Mark
	err = json.Unmarshal(file, &marks)
	if err != nil {
		return nil, err
	}

	return marks, nil
}

func LoadAeroWindowsFixture(jsonFilePath string) ([]aerospaceipc.Window, error) {
	file, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	var windows []aerospaceipc.Window
	err = json.Unmarshal(file, &windows)
	if err != nil {
		return nil, err
	}

	return windows, nil
}

func LoadAeroWindowsFixtureRaw(jsonFilePath string) (string, error) {
	file, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
