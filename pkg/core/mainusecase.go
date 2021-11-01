package core

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type MainUseCaseImpl struct {
	log *log.Logger
}

func NewMainUseCaseImpl(l *log.Logger) UseCase {
	return &MainUseCaseImpl{
		log: l,
	}
}

// RetrieveSecret pulls a single secret from a vault and puts it into a Repository
func (m *MainUseCaseImpl) RetrieveSecret(ctx context.Context, factory Factory, defaults *Defaults, repository Repository, vault *Vault, secret *Secret) error {

	va := factory.NewVaultAccessor(vault.Spec.Type)
	if va == nil {
		return errors.New("internal error: unable to handle vault of given type")
	}

	updatedSecret, err := va.RetrieveSecret(ctx, defaults, vault, secret)
	if err != nil {
		return err
	}

	repository.Put(secret.Name, updatedSecret)

	return nil
}

// Transform applies transformation steps to repository
func (m *MainUseCaseImpl) Transform(ctx context.Context, factory Factory,
	defaults *Defaults, repository Repository, secret *Secret, transformation *Transformations) error {
	return nil
}

// WriteToSink writes output a single sink by pulling it from the repository
func (m *MainUseCaseImpl) WriteToSink(ctx context.Context, factory Factory, defaults *Defaults, repository Repository, sink *Sink) error {

	// get secret to be written from repository
	repositoryContent, err := repository.Get(sink.Var)
	if err != nil {
		return err
	}

	var secret *Secret = repositoryContent.(*Secret)

	// get sik writer for type
	sw := factory.NewSinkWriter(sink.Type)
	if sw == nil {
		return errors.New("internal error: unable to handle sink of given type")
	}

	// write to sink and be done
	err = sw.Write(ctx,defaults,secret,sink)
	if err != nil {
		return err
	}

	return nil
}

func (m *MainUseCaseImpl) Process(ctx context.Context, factory Factory, defaults *Defaults,
	vaults *Vaults, secrets *Secrets, transformations *Transformations, sinks *Sinks) error {

	// need at least one secret, from one vault going to one sink. If either is missing, we cannot proceed.
	if secrets == nil || len(*secrets) == 0 {
		return nil
	}
	if sinks == nil || len(*sinks) == 0 {
		return nil
	}
	if vaults == nil || len(*vaults) == 0 {
		return nil
	}

	repo := factory.NewRepository()

	m.log.Printf("Pulling secrets from vaults")
	for _, secret := range *secrets	{
		vault := vaults.GetVaultByName(secret.VaultName)
		if vault == nil {
			return errors.New(fmt.Sprintf("No such vault: %s", secret.VaultName))
		}
		if err := m.RetrieveSecret(ctx, factory, defaults, repo, vault, secret); err != nil {
			return err
		}
	}

	// Applying transformations
	m.log.Printf("Applying transformations")
	if transformations != nil {
		for _, secret := range *secrets	{
			if err := m.Transform(ctx, factory, defaults, repo, secret, transformations); err != nil {
				return err
			}
		}
	}

	// writing to all sinks
	m.log.Printf("Writing secrets to sinks")
	if sinks != nil {
		for _, sink := range *sinks {
			if err := m.WriteToSink(ctx, factory, defaults, repo, sink); err != nil {
				return err
			}
		}
	}

	return nil
}
