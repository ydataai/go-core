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

// CreateDirectory creates a new directory in the give path
// fails if the the path exists and is not a folder or other reason os related
func (sc *StorageClient) CreateDirectory(relativePath string) error {
	sc.logger.Infof("attempting to create directory %s", relativePath)

	fullPath := fmt.Sprintf("%s/%s", sc.configuration.BasePath, relativePath)
	pathInfo, err := os.Stat(fullPath)

	if os.IsNotExist(err) {
		return os.MkdirAll(fullPath, os.ModePerm)
	}

	if os.IsExist(err) && pathInfo.IsDir() {
		return nil
	}

	return err
}

// RemoveDirectory attempts to remove the directory that holds requirements.txt
func (sc *StorageClient) RemoveDirectory(relativePath string) error {
	fullPath := fmt.Sprintf("%s/%s", sc.configuration.BasePath, relativePath)
	sc.logger.Infof("attempting to remove %s", fullPath)
	if err := os.RemoveAll(fullPath); err != nil {
		sc.logger.Errorf("while remove path %s. Error:", relativePath, err)
		return err
	}
	return nil
}

// CheckIfExists attempts to check if the directory exists
func (sc *StorageClient) CheckIfExists(relativePath string) bool {
	fullPath := fmt.Sprintf("%s/%s", sc.configuration.BasePath, relativePath)
	if _, err := os.Stat(fullPath); err != nil && os.IsNotExist(err) {
		sc.logger.Errorf("while check path %s. Error:", relativePath, err)
		return false
	}

	return true
}

// Rename attempts to rename a path
func (sc *StorageClient) Rename(fromRelativePath, toAbsolutePath string) error {
	sc.logger.Infof("attempting to rename from %s to %s", fromRelativePath, toAbsolutePath)
	if err := os.Rename(fromRelativePath, toAbsolutePath); err != nil {
		sc.logger.Errorf("while rename path from %s to %s. Error:", fromRelativePath, toAbsolutePath, err)
		return err
	}

	return nil
}
