# go-secretshelper

`go-secretshelper` is both a library and a CLI to access secrets stored in vaults such as Cloud-based secrets managers, transform them and store them in files or templates.

[![Go](https://github.com/aschmidt75/go-secretshelper/actions/workflows/go.yaml/badge.svg)](https://github.com/aschmidt75/go-secretshelper/actions/workflows/go.yaml)

## Usage

`go-secretshelper` expects a yaml-based configuration file, which it processes. The configuration contains three
major elements:

* **Vaults** specify, where secrets are stored. Examples are Azure Key Vault or AWS Secrets Manager
* **Transformation** describe, how secrets are modified, e.g. to decode base64 or extract a key from a PKCS#12 container
* **Sinks** specify where and how secrets are written. At present, only files are supported as sinks-

To run a configuration, use: 

```bash
$ go-secretshelper run -c <config file>
```

Sample configuration file:
```yaml
vaults:
  - name: myvault
    type: age-file
    spec:
      path: /some-age-encrypted-file
      identity: /age--private-key

secrets:
  - type: secret
    vault: myvault
    name: sample

sinks:
  - type: file
    var: sample
    spec:
      path: ./sample.dat
      mode: 400
```

See [docs/](docs/README.md) for more details

## Building

The Makefile's `build` target builds an executable in `dist/`.

```bash
$ make build 
```

To build exectuables for several platforms, the `release` target uses [goreleaser](https://goreleaser.com/):

```bash
$ make release
```


## Testing

### Unit tests

```bash
$ go test -v ./...
```

### CLI tests

CLI tests are shell-based and written using bats. The executable is expected to be present in `dist/`. so `make build` 
is necessary before. To run the tests:

```bash
$ cd tests
$ bats .
```

(C) 2021 @aschmidt75, MIT License
