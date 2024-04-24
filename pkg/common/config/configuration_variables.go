package config

// ConfigurationVariables is an interface with LoadFromEnvVars method. The method is meant to be implemented as a struct
// method:
//
//	type StructWithVars struct {
//		Var1 string        `envconfig:"VAR1"`
//		Var2 string        `envconfig:"VAR2"`
//	}
//
//	func (swv *StructWithVars) LoadFromEnvVars() error {
//	  if err := envconfig.Process("", swv); err != nil {
//			return err
//		}
//		return nil
//	}
type ConfigurationVariables interface {
	LoadFromEnvVars() error
}

// InitConfigurationVariables according to environment
func InitConfigurationVariables(configs []ConfigurationVariables) error {
	for _, configuration := range configs {
		if err := configuration.LoadFromEnvVars(); err != nil {
			return err
		}
	}
	return nil
}
