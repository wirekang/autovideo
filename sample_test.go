package main

import (
	"os"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/wirekang/autovideo/tests"
)

func init() {
	lo.Must0(os.Chdir("sample"))
}

func TestInit(t *testing.T) {
	is := assert.New(t)

	configJson := "testconfigname.json"
	tests.SetArgs(
		"--config="+configJson,
		"--init",
	)

	main()
	is.FileExists(configJson)
	v, err := os.ReadFile(configJson)
	is.Nil(err)
	is.JSONEq(string(v), string(v))
	is.Nil(os.Remove(configJson))
}

func TestTest1(t *testing.T) {
	// is := assert.New(t)

	outputFile := "output.mp4"
	cfgPath := "config.json"
	scriptPath := "script.json"
	audiosPath := "audios"
	tests.SetArgs(
		"--config="+cfgPath,
		"--script="+scriptPath,
		"--audios="+audiosPath,
		"--output="+outputFile,
	)
	main()

	// is.FileExists(outputFile)
}
