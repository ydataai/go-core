// Package errors has Fabric Platform base error module.
package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

// FabricError represents a shared error model.
type FabricError interface {
	error
	GetName() string
	GetDescription() string
	GetHTTPCode() int
	GetReturnValue() int
	GetContext() Context
}

// Context is a key/value map to carrier any additional information.
type Context map[string]string

type fabricError struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	HTTPCode    int     `json:"httpCode"`
	ReturnValue int     `json:"returnValue"`
	Context     Context `json:"context,omitempty"`
}

// GetName returns the error name.
func (e fabricError) GetName() string {
	return e.Name
}

// GetDescription returns the error description.
func (e fabricError) GetDescription() string {
	return e.Description
}

// GetHTTPCode returns the HTTP status code of error.
func (e fabricError) GetHTTPCode() int {
	return e.HTTPCode
}

// GetReturnValue returns the exit code of error.
func (e fabricError) GetReturnValue() int {
	return e.ReturnValue
}

// GetContext returns the context of error.
func (e fabricError) GetContext() Context {
	return e.Context
}

// Error interface implementation.
func (e fabricError) Error() string {
	return e.String()
}

// String returns an string representation of BaseError.
func (e fabricError) String() string {
	return fmt.Sprintf("%s (%d) %s", e.Name, e.ReturnValue, e.Description)
}

// ToJSON encode FabircError to JSON string.
func (e fabricError) ToJSON() (string, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(e); err != nil {
		return "<nil>", err
	}
	return buf.String(), nil
}

// New Creates a new BaseError with the required fields.
func New(returnValue, httpCode int, name, description string) FabricError {
	return &fabricError{
		ReturnValue: returnValue,
		Name:        name,
		Description: description,
		HTTPCode:    httpCode,
	}
}

// NewFromJSON creates a Frabric error based on
// JSON string error representation.
func NewFromJSON(text string) FabricError {
	buf := bytes.NewBufferString(text)
	ferr := fabricError{}
	if err := json.NewDecoder(buf).Decode(&ferr); err != nil {
		return &fabricError{
			Name:        "UnexpectedError",
			Description: err.Error(),
			ReturnValue: -1,
			HTTPCode:    500,
		}
	}
	return &ferr
}

// NewFromPod creates a FabricError based on LastTerminationState of ContainerStatus
// given a container Name.
//
// This method is able to decode ContainerStatus.LastTerminationState.Terminated.Message as
// a JSON payload with the following structure:
//
//	{
//	  "name":"UnknownErrorDuringTraining",
//	  "description":"Some unknown and specific error during Synth training either training",
//	  "httpCode":500,
//	  "returnValue":-1,
//	  "context":{
//	    "key1":"value1",
//	    "key2":"value2"
//	  }
//	}
//
// Returns a FabricError or nil if the container is not found.
func NewFromPod(pod corev1.Pod, containerName string) FabricError {
	for _, status := range pod.Status.ContainerStatuses {
		terminated := status.State.Terminated
		if status.Name == containerName && terminated != nil {
			return NewFromJSON(terminated.Message)
		}
	}
	return &fabricError{
		Name:        "NotFoundError",
		Description: fmt.Sprintf("Container %s with Terminated state not found", containerName),
		ReturnValue: -404,
		HTTPCode:    404,
	}
}
