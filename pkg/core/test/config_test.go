package test

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"go-secretshelper/pkg/core"
	"testing"
)

func TestConfig(t *testing.T) {
	if core.NewDefaultConfig() == nil {
		t.Error("Must provide a default config")
	}

	cfg, err := core.NewConfigFromFile("no.such.config.json")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	cfg, err = core.NewConfigFromFile("../../../suppl/fixtures/fixture-1.yaml")
	if err != nil {
		t.Errorf("Expected err=null, got err=%s", err)
	}
	if cfg == nil {
		t.Error("Expected config result, got nil")
	}
}

func DumpValidationErrors(err error) {
	if err != nil {

		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		return
	}
}

func TestValidation(t *testing.T) {
	cfg, err := core.NewConfigFromFile("../../../suppl/fixtures/fixture-1.yaml")
	if err != nil {
		t.Errorf("Expected err=null, got err=%s", err)
	}
	if cfg == nil {
		t.Error("Expected config result, got nil")
	}

	err = cfg.Validate()
	if err != nil {
		t.Errorf("Expected nil got err=%#v", err)
		DumpValidationErrors(err)
	}

	// simple validation (missing elements)

	cfg = &core.Config{
		Vaults: []*core.Vault{
			&core.Vault{
			},
		},
	}
	err = cfg.Validate()
	if err == nil {
		t.Errorf("Expected validation error, got nil")
	}

	// referential validation
	cfg = &core.Config{
		Vaults: []*core.Vault{
			&core.Vault{
				Name: "a",
				Spec: core.VaultSpec{
					Type: "nonex",
					URL: "http://www.example.com",
				},
			},
		},
		Secrets: []*core.Secret{
			&core.Secret{
				Name:      "b",
				VaultName: "nonex",				// this vault is not defined above
				Type:      "secret",
			},
		},
	}
	err = cfg.Validate()
	if err == nil {
		t.Errorf("Expected validation error, got nil")
	}

}
