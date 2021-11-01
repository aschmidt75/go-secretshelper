package core

import "fmt"

type Secrets []*Secret

type Secret struct {
	Name      string `yaml:"name" validate:"required"`
	VaultName string `yaml:"vault" validate:"required"`
	Type      string `yaml:"type" validate:"required"`

	RawContent []byte
}

func (s Secret) String() string {
	set := false
	if len(s.RawContent) > 0 {
		set = true
	}
	return fmt.Sprintf("Secret:[name=%s, Type=%s, set=%t]", s.Name, s.Type, set)
}