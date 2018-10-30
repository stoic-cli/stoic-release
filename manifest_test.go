package release_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stoic-cli/stoic-release"
	"github.com/stoic-cli/stoic-release/mock"
	"github.com/stretchr/testify/assert"
)

func TestManifest(t *testing.T) {
	p := "MyProject"
	v := "v1.0.0"
	s := mock.ValidSignee()
	a := mock.ValidArtifacts()

	m := release.NewManifest(p, v, s, a)

	reader, err := m.Serialise()
	assert.Nil(t, err)

	m2, err := release.NewManifestLoader().Read(reader)
	assert.Nil(t, err)

	r1, err := m2.Serialise()
	assert.Nil(t, err)
	r2, err := m.Serialise()
	assert.Nil(t, err)

	var buf1, buf2 bytes.Buffer
	_, err = io.Copy(&buf1, r1)
	assert.Nil(t, err)
	_, err = io.Copy(&buf2, r2)
	assert.Nil(t, err)
	assert.Equal(t, buf1.String(), buf2.String())
}
