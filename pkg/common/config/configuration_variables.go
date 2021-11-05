package common

import (
	"fmt"
	"os"
	"strconv"
)

// ConfigurationVariables is an interface with LoadFromEnvVars method. The method is meant to be implemented as a struct
// method:
//
//	type StructWithVars struct {
//		var1 string
//		var2 string
//	}
//
//	func (swv *StructWithVars) LoadFromEnvVars() error {
//		swv.var1 = os.Getenv("VALUE_1")
//		swv.var2 = "value 2"
//	}
//
type ConfigurationVariables interface {
	LoadFromEnvVars() error
}

// VariableFromEnvironment check from the value in the os.Getenv and returns it of error if doesn't exists.
func VariableFromEnvironment(env string) (string, error) {
	value := os.Getenv(env)
	if value == "" {
		return "", fmt.Errorf("%s variable not set", env)
	}

	return value, nil
}

// BooleanVariableFromEnvironment check from the value in the os.Getenv, converts to boolean returns it of error if doesn't exists.
func BooleanVariableFromEnvironment(env string) (bool, error) {
	value, err := VariableFromEnvironment(env)
	if err != nil {
		return false, err
	}

	booleanValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, err
	}

	return booleanValue, nil
}

// IntVariableFromEnvironment check from the value in the os.Getenv and returns it of error if doesn't exists.
func IntVariableFromEnvironment(env string) (int, error) {
	value, err := VariableFromEnvironment(env)
	if err != nil {
		return 0, err
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return intValue, nil
}

// Int32VariableFromEnvironment check from the value in the os.Getenv and returns it of error if doesn't exists.
func Int32VariableFromEnvironment(env string) (int32, error) {
	value, err := IntVariableFromEnvironment(env)
	if err != nil {
		return 0, err
	}

	return int32(value), nil
}

// Int64VariableFromEnvironment check from the value in the os.Getenv and returns it of error if doesn't exists.
func Int64VariableFromEnvironment(env string) (int64, error) {
	value, err := IntVariableFromEnvironment(env)
	if err != nil {
		return 0, err
	}

	return int64(value), nil
}
