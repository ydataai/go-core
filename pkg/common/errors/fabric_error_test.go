package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Custom error for testing purpose.

type UnknownErrorDuringTraining struct {
	*FabricError
}

func NewUnknownErrorDuringTraining() error {
	return &UnknownErrorDuringTraining{
		&FabricError{
			Name:        "UnknownErrorDuringTraining",
			Description: "Some unknown and specific error during Synth training either training",
			HTTPCode:    500,
			ReturnValue: -1,
		},
	}
}

// Test functions

func TestUnknownErrorDuringTraining(t *testing.T) {
	ferr := NewUnknownErrorDuringTraining()
	str := "UnknownErrorDuringTraining (-1) Some unknown and specific error during Synth training either training"
	assert.Equal(t, str, ferr.Error())

	expected := "{\"name\":\"UnknownErrorDuringTraining\",\"description\":\"Some unknown and specific error during Synth training either training\",\"httpCode\":500,\"returnValue\":-1}\n"
	actual, err := ferr.(UnknownErrorDuringTraining).ToJSON()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFabricErrorToJSON(t *testing.T) {
	expected := "{\"name\":\"UnknownErrorDuringTraining\",\"description\":\"Some unknown and specific error during Synth training either training\",\"httpCode\":500,\"returnValue\":-1}\n"

	ferr := FabricError{
		Name:        "UnknownErrorDuringTraining",
		Description: "Some unknown and specific error during Synth training either training",
		HTTPCode:    500,
		ReturnValue: -1,
	}
	actual, err := ferr.ToJSON()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFabricErrorFromJSON(t *testing.T) {
	expected := FabricError{
		Name:        "UnknownErrorDuringTraining",
		Description: "Some unknown and specific error during Synth training either training",
		HTTPCode:    500,
		ReturnValue: -1,
	}

	json := "{\"name\":\"UnknownErrorDuringTraining\",\"description\":\"Some unknown and specific error during Synth training either training\",\"httpCode\":500,\"returnValue\":-1}\n"
	actual := NewFromJSON(json)

	assert.EqualValues(t, expected, actual)
}

func TestNewFabricErrorJSON(t *testing.T) {
	expected := "{\"name\":\"UnknownErrorDuringTraining\",\"description\":\"Some unknown and specific error during Synth training either training\",\"returnValue\":-1}\n"

	ferr := New(-1, "UnknownErrorDuringTraining", "Some unknown and specific error during Synth training either training")
	actual, err := ferr.ToJSON()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestErrorString(t *testing.T) {
	err := New(-1, "UnknownErrorDuringTraining", "Some unknown and specific error during Synth training either training")
	str := "UnknownErrorDuringTraining (-1) Some unknown and specific error during Synth training either training"
	assert.Equal(t, str, err.String())
}

func TestErrorError(t *testing.T) {
	err := New(-1, "UnknownErrorDuringTraining", "Some unknown and specific error during Synth training either training")
	str := "UnknownErrorDuringTraining (-1) Some unknown and specific error during Synth training either training"
	assert.Equal(t, str, err.Error())
}
