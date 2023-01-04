package common

type Environment string

const (
	// EnvironmentDevelopment represents when we are working on a development environment
	EnvironmentDevelopment Environment = "dev"
	// EnvironmentProduction represents when we are working on a production environment
	EnvironmentProduction Environment = "prod"
)
