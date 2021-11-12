## Vaults

### AGE files

[age](https://github.com/FiloSottile/age) as an encryption tool can be used as a source for
secrets, e.g. for development purposes. `go-secretshelper` accepts armored, age-encrypted files 
and a matching identity, also from a file. Example:

```yaml
vaults:
- name: kv
  type: age-file
  spec:
    path: ./fixtures/test-agefile
    identity: ./fixtures/test-identity
```

This will retrieve the secrets from the age file under `spec/path` and decrypt them with the identity (from `spec/identity`).
The file has to be json-encoded, mapping the names of secrets to the secrets, e.g. to produce:

```bash
$ age-keygen -o ./fixtures/test-identity
$ echo '{ "test": "s3cr3t" }' | age -e -r <identity-from-previous-step> -a
```

### Azure Key Vault

Secrets can be accessed from an [Azure Key Vault](https://azure.microsoft.com/de-de/services/key-vault/).
Within the `vault` section of a configuration file, add the following to access a vault under
a given URL:

```yaml
  - name: kv
    type: azure-key-vault
    spec:
      url: https://my-sample-vault.vault.azure.net/
```

This will access secrets in `my-sample-vault`. If you want to access secrets in a different
type of vault (e.g. HSM-backed) you can specify the URL accordingly.

In case of default vault service, the url can be omitted. The following snippet does the same.
However this requires using the name of the vault as it is in the `sinks` and `transformations` sections

```yaml
  - name: my-sample-vault
    type: azure-key-vault
```
