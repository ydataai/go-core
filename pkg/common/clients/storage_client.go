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

// NewStorageClient returns an initialized struct with the required dependencies injected
func NewStorageClient(logger *logrus.Logger, configuration StorageClientConfiguration) StorageClient {
	return StorageClient{
		logger:        logger,
		configuration: configuration,
	}
}

// CreateWorkspace attempts to create a directory to hold requirements.txt
func (sc StorageClient) CreateWorkspace(obj, relativePath string) error {
	fullPath := fmt.Sprintf("%s/%s/%s", sc.configuration.BasePath, obj, relativePath)

	sc.logger.Info("Attempting to create directory ", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return os.MkdirAll(fullPath, os.ModePerm)
	}

	return nil
}

// CheckIfExists attempts to check if the directory exists
func (sc *StorageClient) CheckIfExists(obj, relativePath string) bool {
	fullPath := fmt.Sprintf("%s/%s/%s", sc.configuration.BasePath, obj, relativePath)

	_, err := os.Stat(fullPath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// RemoveDirectory attempts to remove the directory that holds requirements.txt
func (sc StorageClient) RemoveDirectory(obj, relativePath string) error {
	fullPath := fmt.Sprintf("%s/%s/%s", sc.configuration.BasePath, obj, relativePath)
	sc.logger.Info("Attempting to remove ", fullPath)
	return os.RemoveAll(fullPath)
}
