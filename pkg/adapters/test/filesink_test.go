package test

import (
	"context"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/adapters"
	"go-secretshelper/pkg/core"
	"log"
	"reflect"
	"testing"
)

func TestFileSink(t *testing.T) {
	secrets := &core.Secrets{
		&core.Secret{
			Name:       "test",
			Type:       "secret",
			VaultName:  "test",
			RawContent: []byte("s3cr3t"),
		},
	}
	sinks := &core.Sinks{
		&core.Sink{
			Type: "mock",
			Path: "test.dat",
			Var: "test",
		},
	}

	fs := afero.NewMemMapFs()

	sink := adapters.NewFileSink(log.Default(), fs)
	err := sink.Write(context.TODO(), &core.Defaults{}, (*secrets)[0], (*sinks)[0])
	if err != nil {
		t.Errorf("Unexpected: %s", err)
	}

	//
	fi, err := fs.Stat((*sinks)[0].Path)
	if err != nil {
		t.Errorf("Unexpected: %s", err)
	}
	if fi.Size() != int64(len((*secrets)[0].RawContent)) {
		t.Errorf("Invalid size")
	}

	raw, err := afero.ReadFile(fs,(*sinks)[0].Path)
	if err != nil {
		t.Errorf("Unexpected: %s", err)
	}

	cmp := reflect.DeepEqual(raw, (*secrets)[0].RawContent)
	if !cmp {
		t.Errorf("Invalid content")
	}
}
