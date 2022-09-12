package clients

import (
	"fmt"
	"os"

	"github.com/ydataai/go-core/pkg/common/logging"
)

// StorageClient is a structure that holds all the dependencies for the following client
type StorageClient struct {
	logger        logging.Logger
	configuration StorageClientConfiguration
}

// StorageClientInterface defines storage client interface
type StorageClientInterface interface {
	BasePath() string
	CreateDirectory(relativePath string) error
	RemoveDirectory(relativePath string) error
	CheckIfExists(relativePath string) bool
	Rename(fromRelativePath, toAbsolutePath string) error
}

// NewStorageClient returns an initialized struct with the required dependencies injected
func NewStorageClient(logger logging.Logger, configuration StorageClientConfiguration) StorageClient {
	return StorageClient{
		logger:        logger,
		configuration: configuration,
	}
}

// BasePath returns a configured base path
func (sc *StorageClient) BasePath() string {
	return sc.configuration.BasePath
}

// CreateDirectory attempts to create a directory to hold requirements.txt
func (sc *StorageClient) CreateDirectory(relativePath string) error {
	fullPath := fmt.Sprintf("%s/%s", sc.configuration.BasePath, relativePath)

	sc.logger.Infof("attempting to create directory %s", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		sc.logger.Errorf("while create path %s", relativePath)
		return os.MkdirAll(fullPath, os.ModePerm)
	}

	return nil
}

// RemoveDirectory attempts to remove the directory that holds requirements.txt
func (sc *StorageClient) RemoveDirectory(relativePath string) error {
	fullPath := fmt.Sprintf("%s/%s", sc.configuration.BasePath, relativePath)
	sc.logger.Infof("attempting to remove %s", fullPath)
	if err := os.RemoveAll(fullPath); err != nil {
		sc.logger.Errorf("while remove path %s", relativePath)
		return err
	}
	return nil
}

// CheckIfExists attempts to check if the directory exists
func (sc *StorageClient) CheckIfExists(relativePath string) error {
	fullPath := fmt.Sprintf("%s/%s", sc.configuration.BasePath, relativePath)
	if _, err := os.Stat(fullPath); err != nil && os.IsNotExist(err) {
		sc.logger.Errorf("while check path %s", relativePath)
		return err
	}

	return nil
}

// Rename attempts to rename a path
func (sc *StorageClient) Rename(fromRelativePath, toAbsolutePath string) bool {
	sc.logger.Infof("attempting to rename from %s to %s", fromRelativePath, toAbsolutePath)
	if err := os.Rename(fromRelativePath, toAbsolutePath); err != nil {
		sc.logger.Errorf("while rename path from %s to %s", fromRelativePath, toAbsolutePath)
		return false
	}

	return true
}
