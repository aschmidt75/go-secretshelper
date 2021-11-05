//go:generate mockgen -package mocks -destination=mocks/mock_usecase.go go-secretshelper/pkg/core UseCase
package core

import "context"

// UseCase is the interface for processing a configuration with all individual steps
type UseCase interface {

	// RetrieveSecret pulls a single secret from a vault and puts it into a Repository
	RetrieveSecret(context.Context, Factory, *Defaults, Repository, *Vault, *Secret) error

	// Transform applies transformation steps to repository
	Transform(context.Context, Factory, *Defaults, Repository, *Secret, *Transformations) error

	// WriteToSink writes output a single sink by pulling it from the repository
	WriteToSink(context.Context, Factory, *Defaults, Repository, *Sink) error

	// Process processes the main use case with given inputs by pulling all secrets,
	// applying transformations and writing to sinks. It pulls port implementations from
	// given Factory.
	Process(context.Context, Factory, *Defaults, *Vaults, *Secrets, *Transformations, *Sinks) error
}
