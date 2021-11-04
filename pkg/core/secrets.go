package core

import "fmt"

// Secrets is an array of Secret structs
type Secrets []*Secret

// Secret defines a named secrets, referenced in a named Vault
type Secret struct {
	Name      string `yaml:"name" validate:"required"`
	VaultName string `yaml:"vault" validate:"required"`
	Type      string `yaml:"type" validate:"required"`

	RawContent []byte
}

// String returns a string representation of a secret
func (s Secret) String() string {
	set := false
	if len(s.RawContent) > 0 {
		set = true
	}
	return fmt.Sprintf("Secret:[name=%s, Type=%s, set=%t]", s.Name, s.Type, set)
}