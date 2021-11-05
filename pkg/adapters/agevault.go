package adapters

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"filippo.io/age"
	"filippo.io/age/armor"
	"fmt"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/core"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
	"strings"
)

const AgePassphrase = "SECRETSHELPER_AGE_PASSPHRASE"

// AgeVault is a core.VaultAccessorPort which pulls secrets from an age-encrypted file
type AgeVault struct {
	log *log.Logger
	fs afero.Fs
}

// NewAgeVault creates a new age vault
func NewAgeVault(log *log.Logger, fs afero.Fs) *AgeVault {
	return &AgeVault{
		log: log,
		fs: fs,
	}
}

type AgeVaultSpec struct {
	Path string
	IdentityFile string
}

func NewAgeVaultSpec(in map[interface{}]interface{}) (AgeVaultSpec, error) {
	var res AgeVaultSpec

	v, ex := in["path"]
	if !ex {
		return res, errors.New("must provide a path element for an age-based vault spec")
	}
	res.Path = v.(string)

	v, ex = in["identity"]
	if !ex {
		return res, errors.New("must provide an identity element for an age-based vault spec")
	}
	res.IdentityFile = v.(string)

	return res, nil
}

// given path to an age file and the identity file, this method decodes
// the file and returns it as a byte array
func (v *AgeVault) readFromAgeFile(path, identity string) ([]byte, error) {
	// parse identity file
	identities, err := v.parseIdentitiesFile(identity)
	if err != nil {
		return nil, err
	}

	var in io.Reader

	inFile, err := v.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()

	in = inFile.(io.Reader)
	b := new(strings.Builder)

	rr := bufio.NewReader(in)
	if start, _ := rr.Peek(len(armor.Header)); string(start) == armor.Header {
		in = armor.NewReader(rr)
	} else {
		in = rr
	}

	r, err := age.Decrypt(in, identities...)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(b, r)
	if err != nil {
		return nil, err
	}

	return []byte(b.String()), nil
}

func (v *AgeVault) RetrieveSecret(ctx context.Context, defaults *core.Defaults,
	vault *core.Vault, secret *core.Secret) (*core.Secret, error) {

	// parse age vault spec
	spec, err := NewAgeVaultSpec(vault.Spec)
	if err != nil {
		return nil, err
	}

	// read file
	res, err := v.readFromAgeFile(spec.Path, spec.IdentityFile)
	if err != nil {
		return nil, err
	}

	// parse json or yaml or as-is
	var data map[string]string
	err = json.Unmarshal(res, &data)
	if err != nil {
		// try yaml
		err = yaml.Unmarshal(res, &data)
		if err != nil {
			// treat the secret as-is
			return &core.Secret{
				RawContent: res,
				RawContentType: "",
				Name: secret.Name,
				Type: secret.Type,
				VaultName: secret.VaultName,
			}, nil
		}
	}
	content, found := data[secret.Name]
	if !found {
		return nil, errors.New(fmt.Sprintf("unable to find secret %s in vault %s", secret.Name, vault.Name))
	}

	return &core.Secret{
		RawContent: []byte(content),
		RawContentType: "",
		Name: secret.Name,
		Type: secret.Type,
		VaultName: secret.VaultName,
	}, nil
}

// parseIdentitiesFile parses a file that contains age or SSH keys. It returns
// one or more of *age.X25519Identity, *agessh.RSAIdentity, *agessh.Ed25519Identity,
// *agessh.EncryptedSSHIdentity, or *EncryptedIdentity.
func (v *AgeVault) parseIdentitiesFile(name string) ([]age.Identity, error) {
	f, err := v.fs.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	b := bufio.NewReader(f)

	// An unencrypted age identity file.
	ids, err := age.ParseIdentities(b)
	if err != nil {
		return nil, fmt.Errorf("failed to read %q: %v", name, err)
	}
	return ids, nil
}
