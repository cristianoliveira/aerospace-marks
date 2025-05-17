package mocks

import (
	"encoding/json"
	"os"

	aerospacecli_mock "github.com/cristianoliveira/aerospace-marks/internal/mocks/aerospacecli"
	storage_mock "github.com/cristianoliveira/aerospace-marks/internal/mocks/storage"
	"github.com/cristianoliveira/aerospace-marks/internal/storage"
	"github.com/cristianoliveira/aerospace-marks/pkgs/aerospacecli"
	"go.uber.org/mock/gomock"
)

// This module contains a set of mock helpers for mocking AeroSpace socket connections
// and storage, easily used in unit tests.

func MockStorageDbClient(ctrl *gomock.Controller) (*storage_mock.MockStorageDbClient, *storage_mock.MockDatabaseConnector) {
	storageDbClient := storage_mock.NewMockStorageDbClient(ctrl)
	storageDbClient.EXPECT().Close().Return(nil).Times(1)

	mockStorage := storage_mock.NewMockDatabaseConnector(ctrl)
	mockStorage.EXPECT().Connect().Return(storageDbClient, nil).Times(1)

	storage.DefaultConnector = mockStorage

	return storageDbClient, mockStorage
}

func MockAerospaceConnection(ctrl *gomock.Controller) (*aerospacecli_mock.MockAeroSpaceSocketConn, *aerospacecli_mock.MockAeroSpaceConnector) {
	mockAeroSpaceConnection := aerospacecli_mock.NewMockAeroSpaceSocketConn(ctrl)
	mockAeroSpaceConnetor := aerospacecli_mock.NewMockAeroSpaceConnector(ctrl)
	mockAeroSpaceConnetor.EXPECT().Connect().Return(mockAeroSpaceConnection, nil).Times(1)
	mockAeroSpaceConnection.EXPECT().CloseConnection().Return(nil).Times(1)

	aerospacecli.DefaultConnector = mockAeroSpaceConnetor

	return mockAeroSpaceConnection, mockAeroSpaceConnetor
}

func MockStorageDbResult(ctrl *gomock.Controller, lastInsertId *int64, rowsAffected *int64) (*storage_mock.MockDbResult) {
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

func LoadAeroWindowsFixture(jsonFilePath string) ([]aerospacecli.Window, error) {
	file, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	var windows []aerospacecli.Window
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
