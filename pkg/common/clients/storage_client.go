package clients

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// StorageClient is a structure that holds all the dependencies for the following client
type StorageClient struct {
	logger        *logrus.Logger
	configuration StorageClientConfiguration
}

// StorageClientInterface defines storage client interface
type StorageClientInterface interface {
	CreateDirectory(relativePath string) error
	RemoveDirectory(relativePath string) error
	CheckIfExists(relativePath string) bool
}

// NewStorageClient returns an initialized struct with the required dependencies injected
func NewStorageClient(logger *logrus.Logger, configuration StorageClientConfiguration) StorageClient {
	return StorageClient{
		logger:        logger,
		configuration: configuration,
	}
}

// CreateDirectory attempts to create a directory to hold requirements.txt
func (sc *StorageClient) CreateDirectory(relativePath string) error {
	fullPath := fmt.Sprintf("%s/%s", sc.configuration.BasePath, relativePath)

	sc.logger.Info("Attempting to create directory ", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return os.MkdirAll(fullPath, os.ModePerm)
	}

	return nil
}

// RemoveDirectory attempts to remove the directory that holds requirements.txt
func (sc *StorageClient) RemoveDirectory(relativePath string) error {
	fullPath := fmt.Sprintf("%s/%s", sc.configuration.BasePath, relativePath)
	sc.logger.Info("Attempting to remove ", fullPath)
	return os.RemoveAll(fullPath)
}

// CheckIfExists attempts to check if the directory exists
func (sc *StorageClient) CheckIfExists(relativePath string) bool {
	fullPath := fmt.Sprintf("%s/%s", sc.configuration.BasePath, relativePath)

	_, err := os.Stat(fullPath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}
