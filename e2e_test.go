package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wirekang/autovideo/tests"
)

func TestInit(t *testing.T) {
	is := assert.New(t)
	tests.CdTempDir()

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
}

func TestTest1(t *testing.T) {
	is := assert.New(t)
	tests.CdTempDir()

	outputFile := "output.mp4"
	cfgPath := "testconfigname.json"
	scriptPath := "script.json"
	tests.SetArgs("--config="+cfgPath, "--output="+outputFile, "--script="+scriptPath)

	is.Nil(os.WriteFile(cfgPath, []byte(fmt.Sprintf(`
{
  "image_width": 1280,
  "image_height": 720,
  "font_color": "#fff",
  "background_color": "#000",
  "font_name": "D2Coding.ttf",
  "font_size": 64
}
`)), 0666))

	is.Nil(os.WriteFile(scriptPath, []byte(`
[
	{ "text": "This is text index 0", "millis": 4000 },
	{ "text": "This is text index 1", "millis": 3001 },
	{ "text": "This is text index 2", "millis": 2002 },
	{ "text": "This is text index 3", "millis": 1003 },
	{ "text": "This is text index 4", "millis": 904 },
	{ "text": "This is text index 5", "millis": 805 },
	{ "text": "This is text index 6", "millis": 706 },
	{ "text": "This is text index 7", "millis": 607 },
	{ "text": "This is text index 8", "millis": 508 },
	{ "text": "This is text index 9", "millis": 409 },
	{ "text": "This is text index 10", "millis": 310 }
]
`), 0666))

	main()

	is.FileExists(outputFile)
}
