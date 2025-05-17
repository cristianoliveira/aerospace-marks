package mocks

import (
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
