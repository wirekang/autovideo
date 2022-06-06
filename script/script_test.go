package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	is := assert.New(t)
	title, lines, err := Parse("../sample/script.txt")
	is.Nil(err)
	is.Equal("This is very very very long long long title line", title)
	is.Len(lines, 5)
}
