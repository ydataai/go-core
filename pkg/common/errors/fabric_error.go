// Package errors has Fabric Platform base error module.
package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// FabricError represents a shared error model.
type FabricError struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTTPCode    int    `json:"httpCode,omitempty"`
	ReturnValue int    `json:"returnValue"`
}

// New Creates a new BaseError.
func New(returnValue int, name, description string) FabricError {
	return FabricError{
		ReturnValue: returnValue,
		Name:        name,
		Description: description,
	}
}

// NewEmpty creates a new empty BaseError.
func NewEmpty() FabricError {
	return FabricError{}
}

// NewFromJSON creates a Frabric error based on
// JSON string error representation.
func NewFromJSON(text string) FabricError {
	buf := bytes.NewBufferString(text)
	ferr := FabricError{}
	if err := json.NewDecoder(buf).Decode(&ferr); err != nil {
		return FabricError{}
	}
	return ferr
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
