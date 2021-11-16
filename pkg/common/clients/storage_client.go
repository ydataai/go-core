package clients

import (
	"fmt"
	"os"

	"github.com/ydataai/go-core/pkg/common/logging"
)

// StorageClient is a structure that holds all the dependencies for the following client
type StorageClient struct {
	logger        logging.Logger
	pathSufix     string
	configuration StorageClientConfiguration
}

// StorageClientInterface defines storage client interface
type StorageClientInterface interface {
	CreateDirectory(relativePath string) error
	RemoveDirectory(relativePath string) error
	CheckIfExists(relativePath string) bool
}

// NewStorageClient returns an initialized struct with the required dependencies injected
func NewStorageClient(logger logging.Logger, pathSufix string, configuration StorageClientConfiguration) StorageClient {
	return StorageClient{
		logger:        logger,
		pathSufix:     pathSufix,
		configuration: configuration,
	}
}

// CreateDirectory attempts to create a directory to hold requirements.txt
func (sc *StorageClient) CreateDirectory(relativePath string) error {
	fullPath := fmt.Sprintf("%s/%s/%s", sc.configuration.BasePath, sc.pathSufix, relativePath)

	sc.logger.Info("Attempting to create directory ", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return os.MkdirAll(fullPath, os.ModePerm)
	}

	return nil
}

// RemoveDirectory attempts to remove the directory that holds requirements.txt
func (sc *StorageClient) RemoveDirectory(relativePath string) error {
	fullPath := fmt.Sprintf("%s/%s/%s", sc.configuration.BasePath, sc.pathSufix, relativePath)
	sc.logger.Info("Attempting to remove ", fullPath)
	return os.RemoveAll(fullPath)
}

// CheckIfExists attempts to check if the directory exists
func (sc *StorageClient) CheckIfExists(relativePath string) bool {
	fullPath := fmt.Sprintf("%s/%s/%s", sc.configuration.BasePath, sc.pathSufix, relativePath)

	_, err := os.Stat(fullPath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}
