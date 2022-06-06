package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wirekang/autovideo/file"
)

func TestMerge(t *testing.T) {
	const output = "output.mp3"
	is := assert.New(t)
	files, err := file.ListByModTime("../sample/audios")
	is.Nil(err)

	is.Nil(Merge(files, output))

	is.FileExists(output)
	file.DeleteOnExit(output)
	file.Delete()
}
