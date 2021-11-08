package core

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

// Config is the main configuration struct given by YAML input
type Config struct {
	Defaults        Defaults        `yaml:"defaults" validate:""`
	Vaults          Vaults          `yaml:"vaults" validate:"required,dive"`
	Secrets         Secrets         `yaml:"secrets" validate:"required,dive"`
	Transformations Transformations `yaml:"transformations" validate:""`
	Sinks           Sinks           `yaml:"sinks" validate:"required,dive"`
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
func (c *Config) IsVarDefined(varName string) bool {
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
func (c *Config) Validate(f Factory) error {
	v := validator.New()

	st := make(map[string]struct{})
	for _, e := range f.SinkTypes() {
		st[e] = struct{}{}
	}
	vat := make(map[string]struct{})
	for _, e := range f.VaultAccessorTypes() {
		vat[e] = struct{}{}
	}

	//
	for _, vault := range c.Vaults {
		if err := v.Struct(vault); err != nil {
			return err
		}

		if _, ex := vat[vault.Type]; !ex {
			return fmt.Errorf("unknown vault type: %s in vault: %s", vault.Type, vault.Name)
		}
	}

	for _, secret := range c.Secrets {
		if err := v.Struct(secret); err != nil {
			return err
		}

		if v := c.Vaults.GetVaultByName(secret.VaultName); v == nil {
			return fmt.Errorf("invalid vault %s referenced in secret %s", secret.VaultName, secret.Name)
		}
	}

	for _, sink := range c.Sinks {
		if err := v.Struct(sink); err != nil {
			return err
		}

		if _, ex := st[sink.Type]; !ex {
			return fmt.Errorf("unknown sink type: %s", sink.Type)
		}

		if !c.IsVarDefined(sink.Var) {
			return fmt.Errorf("invalid variable %s referenced in a sink", sink.Var)
		}
	}

	return nil
}
