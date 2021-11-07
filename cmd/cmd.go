package main

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/adapters"
	"go-secretshelper/pkg/core"
	"log"
	"os"
)

func main() {

	log := createLogger()
	f := adapters.NewBuiltinFactory(log, afero.NewOsFs())
	cmd := core.NewMainUseCaseImpl(log)

	config := core.NewDefaultConfig()

	err := cmd.Process(context.Background(), f,
		&config.Defaults,
		&config.Vaults,
		&config.Secrets,
		&config.Transformations,
		&config.Sinks)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func createLogger() *log.Logger {
	return log.Default()
}