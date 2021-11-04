package test

import (
	"context"
	"github.com/golang/mock/gomock"
	"go-secretshelper/pkg/core"
	"log"
	"testing"
)

func TestMainUseCase(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ctx := context.TODO()

	mf := NewMockFactory(mockCtrl, t)

	vaults := &core.Vaults{
		&core.Vault{
			Name: "test",
			Type: "mock",
			Spec: core.VaultSpec{},
		},
	}
	secrets := &core.Secrets{
		&core.Secret{
			Name: "test",
			Type: "secret",
			VaultName: "test",
		},
	}
	transformations := &core.Transformations{}
	sinks := &core.Sinks{
		&core.Sink{
			Type: "mock",
			Path: "test",
			Var: "test",
		},
	}
	defaults := &core.Defaults{}

	useCase := core.NewMainUseCaseImpl(log.Default())

	// set up expectations
	mf.GetMockVaultAccessor("mock").EXPECT().RetrieveSecret(ctx, defaults, (*vaults)[0], (*secrets)[0]).Return((*secrets)[0], nil).Times(1)
	mf.GetMockRepository().EXPECT().Put("test", (*secrets)[0]).Times(1)
	mf.GetMockRepository().EXPECT().Get("test").Return((*secrets)[0], nil).Times(1)
	mf.GetMockSinkWriter("mock").EXPECT().Write(ctx, defaults, (*secrets)[0], (*sinks)[0]).Times(1)

	err := useCase.Process(ctx, mf, defaults, vaults, secrets, transformations, sinks)
	if err != nil {
		t.Errorf("Unexpected: %s", err)
	}


}
