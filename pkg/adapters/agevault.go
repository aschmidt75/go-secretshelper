package adapters

import (
	"context"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/core"
	"log"
)

// AgeVault is a core.VaultAccessorPort which pulls secrets from an age-encrypted file
type AgeVault struct {
	log *log.Logger
	fs afero.Fs
}

func NewAgeVault(log *log.Logger, fs afero.Fs) *AgeVault {
	return &AgeVault{
		log: log,
		fs: fs,
	}
}

func (v *AgeVault) RetrieveSecret(ctx context.Context, defaults *core.Defaults,
	vault *core.Vault, secret *core.Secret) (*core.Secret, error) {

	return secret, nil
}
