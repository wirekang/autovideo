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

func TestImageSave(t *testing.T) {
	is := assert.New(t)
	tests.CdTempDir()

	outputDir := "outputDir"
	cfgPath := "testconfigname.json"
	linesPath := "lines.json"
	tests.SetArgs("--config="+cfgPath, "--output="+outputDir, linesPath)
	is.Nil(os.WriteFile(cfgPath, []byte(fmt.Sprintf(`
{
  "image_width": 1280,
  "image_height": 720,
  "font_color": "#000",
  "background_color": "#fff",
  "font_name": "D2Coding.ttf",
  "file_prefix": "image",
  "font_size": 64
}
`)), 0666))

	is.Nil(os.WriteFile(linesPath, []byte(`
[
	{ "text": "This is text index 0", "millis": 1000 },
	{ "text": "This is text index 1", "millis": 1001 },
	{ "text": "This is text index 2", "millis": 1002 },
	{ "text": "This is text index 3", "millis": 1003 },
	{ "text": "This is text index 4", "millis": 1004 },
	{ "text": "This is text index 5", "millis": 1005 },
	{ "text": "This is text index 6", "millis": 1006 },
	{ "text": "This is text index 7", "millis": 1007 },
	{ "text": "This is text index 8", "millis": 1008 },
	{ "text": "This is text index 9", "millis": 1009 },
	{ "text": "This is text index 10", "millis": 1010 }
]
`), 0666))

	main()

}
