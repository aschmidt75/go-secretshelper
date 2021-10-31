package test

import (
	"github.com/golang/mock/gomock"
	"go-secretshelper/pkg/core"
	"go-secretshelper/pkg/core/mocks"
	"testing"
)

// MockFactory produces mocks only
type MockFactory struct {
	mockCtrl *gomock.Controller
	t *testing.T

	repo *mocks.MockRepository
	vaults map[string]*mocks.MockVaultAccessorPort
	sinks map[string]*mocks.MockSinkWriterPort
}

func NewMockFactory(mockCtrl *gomock.Controller, t *testing.T) *MockFactory {
	mf := &MockFactory{
		mockCtrl: mockCtrl,
		t: t,
		vaults: make(map[string]*mocks.MockVaultAccessorPort),
		sinks: make(map[string]*mocks.MockSinkWriterPort),
		repo: mocks.NewMockRepository(mockCtrl),
	}

	// auto set up mock port
	mf.newVaultAccessorInternal("mock")
	mf.newSinkWriterInternal("mock")

	return mf
}

func (df *MockFactory) SinkTypes() []string {
	return []string{
		"mock",
	}
}

func (df *MockFactory) TransformationTypes() []string {
	return []string{
		"mock",
	}
}

func (df *MockFactory) VaultAccessorTypes() []string {
	return []string{
		"mock",
	}
}

func (df *MockFactory) NewRepository() core.Repository {
	return df.repo
}

func (df *MockFactory) GetMockRepository() *mocks.MockRepository {
	return df.repo

}

func (df *MockFactory) newSinkWriterInternal(sinkType string) core.SinkWriterPort {
	s := mocks.NewMockSinkWriterPort(df.mockCtrl)
	df.sinks[sinkType] = s
	return s
}

func (df *MockFactory) GetMockSinkWriter(sinkType string) *mocks.MockSinkWriterPort {
	return df.sinks[sinkType]
}

func (df *MockFactory) NewSinkWriter(sinkType string) core.SinkWriterPort {
	return df.sinks[sinkType]
}

func (df *MockFactory) NewTransformation(transformationType string) core.TransformationPort {
	return mocks.NewMockTransformationPort(df.mockCtrl)
}

func (df *MockFactory) NewVaultAccessor(vaultType string) core.VaultAccessorPort {
	return df.vaults[vaultType]
}


func (df *MockFactory) newVaultAccessorInternal(vaultType string) core.VaultAccessorPort {
	df.t.Logf("MockFactory.NewVaultAccessor, type=%s\n", vaultType)
	va := mocks.NewMockVaultAccessorPort(df.mockCtrl)
	df.vaults[vaultType] = va
	return va
}

func (df *MockFactory) GetMockVaultAccessor(t string) *mocks.MockVaultAccessorPort {
	return df.vaults[t]
}


