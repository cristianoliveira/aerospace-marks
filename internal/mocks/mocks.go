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
	"github.com/cristianoliveira/aerospace-marks/internal/storage/db/queries"
	"go.uber.org/mock/gomock"
)

// This module contains a set of mock helpers for mocking AeroSpace socket connections
// and storage, easily used in unit tests.

func MockStorageDbClient(ctrl *gomock.Controller) (
	*storage_mock.MockStorageDbClient,
	*storage_mock.MockMarkStorage,
) {
	storageDbClient := storage_mock.NewMockStorageDbClient(ctrl)
	newStorage := storage_mock.NewMockMarkStorage(ctrl)

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

func LoadMarksFixture(jsonFilePath string) ([]queries.Mark, error) {
	file, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	var marks []queries.Mark
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
