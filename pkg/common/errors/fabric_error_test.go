package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			Context: Context{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
}

// Test functions

func TestUnknownErrorDuringTraining(t *testing.T) {
	ferr := NewUnknownErrorDuringTraining()
	str := "UnknownErrorDuringTraining (-1) Some unknown and specific error during Synth training either training"
	assert.Equal(t, str, ferr.Error())

	expected := "{\"name\":\"UnknownErrorDuringTraining\",\"description\":\"Some unknown and specific error during Synth training either training\",\"httpCode\":500,\"returnValue\":-1,\"context\":{\"key1\":\"value1\",\"key2\":\"value2\"}}\n"
	actual, err := ferr.(*UnknownErrorDuringTraining).ToJSON()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFabricErrorWithContextToJSON(t *testing.T) {
	expected := "{\"name\":\"UnknownErrorDuringTraining\",\"description\":\"Some unknown and specific error during Synth training either training\",\"httpCode\":500,\"returnValue\":-1,\"context\":{\"key1\":\"value1\",\"key2\":\"value2\"}}\n"

	ferr := FabricError{
		Name:        "UnknownErrorDuringTraining",
		Description: "Some unknown and specific error during Synth training either training",
		HTTPCode:    500,
		ReturnValue: -1,
		Context: Context{
			"key1": "value1",
			"key2": "value2",
		},
	}
	actual, err := ferr.ToJSON()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestFabricErrorToJSON(t *testing.T) {
	expected := "{\"name\":\"UnknownErrorDuringTraining\",\"description\":\"Some unknown and specific error during Synth training either training\",\"httpCode\":500,\"returnValue\":-1}\n"

	ferr := New(-1, 500, "UnknownErrorDuringTraining", "Some unknown and specific error during Synth training either training")
	actual, err := ferr.ToJSON()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestNewFromJSON(t *testing.T) {
	expected := &FabricError{
		Name:        "UnknownErrorDuringTraining",
		Description: "Some unknown and specific error during Synth training either training",
		HTTPCode:    500,
		ReturnValue: -1,
	}

	json := "{\"name\":\"UnknownErrorDuringTraining\",\"description\":\"Some unknown and specific error during Synth training either training\",\"httpCode\":500,\"returnValue\":-1}\n"
	actual := NewFromJSON(json)

	assert.EqualValues(t, expected, actual)
}

func TestNewFromJSONWithContextNull(t *testing.T) {
	expected := &FabricError{
		Name:        "UnknownErrorDuringTraining",
		Description: "Some unknown and specific error during Synth training either training",
		HTTPCode:    500,
		ReturnValue: -1,
	}

	json := "{\"name\":\"UnknownErrorDuringTraining\",\"description\":\"Some unknown and specific error during Synth training either training\",\"httpCode\":500,\"returnValue\":-1,\"context\":null}\n"
	actual := NewFromJSON(json)

	assert.EqualValues(t, expected, actual)
}

func TestNewFromTerminatedPod(t *testing.T) {
	pod := corev1.Pod{
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name: "sidecar",
					State: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							ExitCode: -1,
							Message:  "Some message",
						},
					},
				},
				{
					Name: "main",
					State: corev1.ContainerState{
						Terminated: &corev1.ContainerStateTerminated{
							ExitCode: -1,
							Message:  "{\"name\":\"UnknownErrorDuringTraining\",\"description\":\"Some unknown and specific error during Synth training either training\",\"httpCode\":500,\"returnValue\":-1}",
						},
					},
				},
			},
		},
	}

	expected := &FabricError{
		Name:        "UnknownErrorDuringTraining",
		Description: "Some unknown and specific error during Synth training either training",
		HTTPCode:    500,
		ReturnValue: -1,
	}

	actual := NewFromPod(pod, "main")

	assert.EqualValues(t, expected, actual)
}

func TestNewFromRunningPod(t *testing.T) {
	pod := corev1.Pod{
		Status: corev1.PodStatus{
			ContainerStatuses: []corev1.ContainerStatus{
				{
					Name: "sidecar",
					State: corev1.ContainerState{
						Running: &corev1.ContainerStateRunning{
							StartedAt: v1.Now(),
						},
					},
				},
				{
					Name: "main",
					State: corev1.ContainerState{
						Running: &corev1.ContainerStateRunning{
							StartedAt: v1.Now(),
						},
					},
				},
			},
		},
	}

	ferr := NewFromPod(pod, "main")
	expected := &FabricError{
		Name:        "NotFoundError",
		Description: "Container main with Terminated state not found",
		ReturnValue: -404,
		HTTPCode:    404,
	}

	assert.EqualValues(t, expected, ferr)
}

func TestErrorString(t *testing.T) {
	err := New(-1, 500, "UnknownErrorDuringTraining", "Some unknown and specific error during Synth training either training")
	str := "UnknownErrorDuringTraining (-1) Some unknown and specific error during Synth training either training"
	assert.Equal(t, str, err.String())
}

func TestErrorError(t *testing.T) {
	err := New(-1, 500, "UnknownErrorDuringTraining", "Some unknown and specific error during Synth training either training")
	str := "UnknownErrorDuringTraining (-1) Some unknown and specific error during Synth training either training"
	assert.Equal(t, str, err.Error())
}
