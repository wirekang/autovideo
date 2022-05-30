package tests

import (
	"os"

	"github.com/spf13/pflag"
)

func CdTempDir() {
	dir, err := os.MkdirTemp("", "temp*")
	if err != nil {
		panic(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func SetArgs(v ...string) {
	os.Args = append([]string{
		os.Args[0],
	}, v...)
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}
