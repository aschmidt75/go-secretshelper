package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spf13/afero"
	"go-secretshelper/pkg/adapters"
	"go-secretshelper/pkg/core"
	"io/ioutil"
	"log"
	"os"
)

var (
	commit  = "none"
	date    = "unknown"
)

func usage() {
	fmt.Println("Usage: go-secretshelper [-v] [-c config] <command>")
	fmt.Println("where commands are")
	fmt.Println("  version		print out version")
	fmt.Println("  run			run specified config")
}

func main() {

	verboseFlag := flag.Bool("v", false, "Enables verbose output")
	flag.Parse()

	var l *log.Logger
	if *verboseFlag {
		l = log.New(os.Stderr, "", log.LstdFlags)
	} else {
		l = log.New(ioutil.Discard, "", 0)
	}

	values := flag.Args()
	if len(values) == 0 {
		usage()
		os.Exit(1)
	}

	switch values[0] {
	case "version":
		fmt.Printf("%s (%s)\n", commit, date)
		os.Exit(0)

	case "run":

		fs := flag.NewFlagSet("run", flag.ExitOnError)
		configFlag := fs.String("c", "", "configuration file")

		if err := fs.Parse(values[1:]); err != nil {
			println(err)
			os.Exit(1)
		}

		// read config
		config, err := core.NewConfigFromFile(*configFlag)
		if err != nil {
			fmt.Printf("Unable to read config from file %s: %s\n", *configFlag, err)
			os.Exit(1)
		}

		// validate
		f := adapters.NewBuiltinFactory(l, afero.NewOsFs())
		if err := config.Validate(f); err != nil {
			fmt.Printf("Error validating configuration: %s\n", err)
			os.Exit(2)
		}

		// run
		cmd := core.NewMainUseCaseImpl(l)

		err = cmd.Process(context.Background(), f,
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

}

func createLogger() *log.Logger {
	return log.Default()
}