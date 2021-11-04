package adapters

import (
	"github.com/spf13/afero"
	"go-secretshelper/pkg/core"
	"go-secretshelper/pkg/core/mocks"
	"log"
)

// BuiltinFactory is able to create all builtin components of
// the adapters package
type BuiltinFactory struct {
	log *log.Logger
	fs afero.Fs
}

// NewBuiltinFactory creates the Builtin Factory
func NewBuiltinFactory(log *log.Logger, fs afero.Fs) *BuiltinFactory {
	return &BuiltinFactory{
		log: log,
		fs: fs,
	}
}

func (f *BuiltinFactory) SinkTypes() []string {
	return []string{
		FileSinkType,
	}
}

func (f *BuiltinFactory) TransformationTypes() []string {
	return []string{}
}

func (f *BuiltinFactory) VaultAccessorTypes() []string {
	return []string{}
}

func (f *BuiltinFactory) NewRepository() core.Repository {
	return nil
}

func (f *BuiltinFactory) GetMockRepository() *mocks.MockRepository {
	return nil

}

func (f *BuiltinFactory) newSinkWriterInternal(sinkType string) core.SinkWriterPort {
	return nil
}

func (f *BuiltinFactory) GetMockSinkWriter(sinkType string) *mocks.MockSinkWriterPort {
	return nil
}

func (f *BuiltinFactory) NewSinkWriter(sinkType string) core.SinkWriterPort {
	switch sinkType {
	case FileSinkType: return NewFileSink(f.log, f.fs)
	}
	return nil
}

func (f *BuiltinFactory) NewTransformation(transformationType string) core.TransformationPort {
	return nil
}

func (f *BuiltinFactory) NewVaultAccessor(vaultType string) core.VaultAccessorPort {
	return nil
}

