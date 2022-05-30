package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wirekang/autovideo/line"
)

func TestParseLines(t *testing.T) {
	is := assert.New(t)
	v := `
[
	{ "text": "string 1", "millis": 1000 },
	{ "text": "string 2", "millis": 1000 },
	{ "text": "string 3", "millis": 1000 },
	{ "text": "string 4", "millis": 1000 },
	{ "text": "string 5", "millis": 1000 }
]
`
	lines, err := line.ParseLines([]byte(v))
	is.Nil(err)
	is.Len(lines, 5)
}

func TestParseLinesTooShort(t *testing.T) {
	is := assert.New(t)
	v := `
[
	{ "text": "string 1", "millis": 1000 },
	{ "text": "string 2", "millis": 1000 }
]
`
	_, err := line.ParseLines([]byte(v))
	is.NotNil(err)

	v = `
[
	{ "text": "string 1", "millis": 0 },
	{ "text": "string 2", "millis": 0 },
	{ "text": "string 2", "millis": 0 },
	{ "text": "string 2", "millis": 1000 },
	{ "text": "string 2", "millis": 0 },
	{ "text": "string 2", "millis": 0 }
]
`
	_, err = line.ParseLines([]byte(v))
	is.NotNil(err)
}
