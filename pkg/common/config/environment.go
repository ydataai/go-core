package common

import (
	"errors"
	"os"
)

// Environment is an enum with the representation of our product environments
type Environment string

const (
	// EnvironmentDev means that the application is running on a development environment. Output may be more verbose
	EnvironmentDev Environment = "development"
	// EnvironmentStaging means that the application is running on a staging environment. Should be similar to production,
	// although it might be unstable
	EnvironmentStaging Environment = "staging"
	// EnvironmentProd means that the application is running on a production environment. Output may be less verbose
	EnvironmentProd Environment = "production"
)

// LoadEnvironmentVar attempts to load the ENVIRONMENT env variable, as it's requested in several places.
// Throws an error if the validations aren't met
func LoadEnvironmentVar() (Environment, error) {
	environment := os.Getenv("ENVIRONMENT")

	if environment == "" {
		return "", errors.New("ENVIRONMENT variable not set")
	}

	var envValue Environment
	switch environment {
	case "dev", "development":
		envValue = EnvironmentDev
	case "staging":
		envValue = EnvironmentStaging
	case "prod", "production":
		envValue = EnvironmentProd
	default:
		return "", errors.New("invalid ENVIRONMENT variable value set")
	}

	return envValue, nil
}
