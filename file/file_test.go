package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMediaLengthSeconds(t *testing.T) {
	const delta = 1
	is := assert.New(t)
	length, err := MediaLengthSeconds("../sample/sample_30s.mp4")
	is.Nil(err)
	is.InDelta(30, length, delta)

	length, err = MediaLengthSeconds("../sample/sample_33s.wav")
	is.Nil(err)
	is.InDelta(33, length, delta)
}
