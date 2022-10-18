// Package errors has Fabric Platform base error module.
package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

type Context map[string]string

// FabricError represents a shared error model.
type FabricError struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	HTTPCode    int     `json:"httpCode,omitempty"`
	ReturnValue int     `json:"returnValue"`
	Context     Context `json:"context,omitempty"`
}

// New Creates a new BaseError with the required fields.
func New(returnValue int, name, description string) FabricError {
	return FabricError{
		ReturnValue: returnValue,
		Name:        name,
		Description: description,
		Context:     Context{},
	}
}

// NewFromJSON creates a Frabric error based on
// JSON string error representation.
func NewFromJSON(text string) *FabricError {
	buf := bytes.NewBufferString(text)
	ferr := FabricError{}
	if err := json.NewDecoder(buf).Decode(&ferr); err != nil {
		return &FabricError{
			Name:        "UnexpectedError",
			Description: err.Error(),
			ReturnValue: -1,
			HTTPCode:    500,
			Context:     Context{},
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
func NewFromPod(pod corev1.Pod, containerName string) *FabricError {
	for _, status := range pod.Status.ContainerStatuses {
		if status.Name == containerName && status.LastTerminationState.Terminated != nil {
			return NewFromJSON(status.LastTerminationState.Terminated.Message)
		}
	}
	return nil
}

// Error interface implementation.
func (e *FabricError) Error() string {
	return e.String()
}

// String returns an string representation of BaseError.
func (e *FabricError) String() string {
	return fmt.Sprintf("%s (%d) %s", e.Name, e.ReturnValue, e.Description)
}

// ToJSON encode FabircError to JSON string.
func (e *FabricError) ToJSON() (string, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(e); err != nil {
		return "<nil>", err
	}
	return buf.String(), nil
}
