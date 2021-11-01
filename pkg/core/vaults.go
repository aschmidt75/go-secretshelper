package core

import (
	"context"
	"fmt"
)

type Vaults []*Vault

type Vault struct {
	Name string `yaml:"name" validate:"required"`
	Spec VaultSpec `yaml:"spec" validate:"required"`
}

type VaultSpec struct {
	Type string `yaml:"type" validate:"required"`
	URL string `yaml:"url" validate:"required,url"`
}

func (v Vault) String() string {
	return fmt.Sprintf("Vault:[Name=%s, Type=%s]", v.Name, v.Spec.Type)
}

func (vaults *Vaults) GetVaultByName(name string) *Vault {
	for _, vault := range *vaults {
		if vault.Name == name {
			return vault
		}
	}
	return nil
}

// VaultAccessorPort is able to pull secrets from a Vault
type VaultAccessorPort interface {

	// RetrieveSecret retrieves a secret from given vault
	RetrieveSecret(context.Context, *Defaults, *Vault, *Secret) (*Secret, error)
}
