package core

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// Config is the main configuration struct given by YAML input
type Config struct {
	Defaults Defaults `yaml:"defaults" validate:""`
	Vaults Vaults `yaml:"vaults" validate:"required,dive"`
	Secrets Secrets `yaml:"secrets" validate:"required,dive"`
	Transformations Transformations `yaml:"transformations" validate:""`
	Sinks Sinks `yaml:"sinks" validate:"required,dive"`
}

// NewDefaultConfig returns a configuration struct with valid default settings
func NewDefaultConfig() *Config {
	return &Config{}
}

// NewConfig is the default way of reading configuration from yaml stream
func NewConfig(in io.Reader) (*Config, error) {
	yamlDec := yaml.NewDecoder(in)

	res := NewDefaultConfig()
	if err := yamlDec.Decode(res); err != nil {
		return res, err
	}

	return res, nil
}

// NewConfigFromFile creates a configuration from yaml file
func NewConfigFromFile(fileName string) (*Config, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return NewConfig(f)
}

// IsVarDefined checks if given variable name is defined, either in
// secrets or as the result of a transformation step
func (c* Config) IsVarDefined(varName string) bool {
	for _, secret := range c.Secrets {
		if secret.Name == varName {
			return true
		}
	}

	for _, transformation := range c.Transformations {
		if transformation.ToVar == varName {
			return true
		}
	}

	return false
}

// Validate validates a configuration using the validator and
// additional cross checks.
// TODO inject Factory for list of valid types.
func (c *Config) Validate() error {
	v := validator.New()

	//
	for _, vault := range c.Vaults {
		if err := v.Struct(vault); err != nil {
			return err
		}
	}

	for _, secret := range c.Secrets {
		if err := v.Struct(secret); err != nil {
			return err
		}

		if v := c.Vaults.GetVaultByName(secret.VaultName); v == nil {
			return errors.New(
				fmt.Sprintf("invalid vault %s referenced in secret %s", secret.VaultName, secret.Name))
		}
	}

	for _, sink := range c.Sinks {
		if err := v.Struct(sink); err != nil {
			return err
		}

		if !c.IsVarDefined(sink.Var) {
			return errors.New(
				fmt.Sprintf("invalid variable %s referenced in a sink", sink.Var))
		}
	}

	return nil
}
