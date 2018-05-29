package release_test

import (
	"testing"
	"github.com/stoic-cli/stoic-release"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"github.com/stoic-cli/stoic-release/mock"
)

func TestNewRelease(t *testing.T) {
	artifacts := mock.ValidArtifacts()
	manifest := release.NewManifest("MyProject", "v1.0.0", mock.ValidSignee(), artifacts)
	signature := []byte("some kind of signature")

	dir, err := ioutil.TempDir("", "release-")
	assert.Nil(t, err)
	saver := release.NewFileSystemSaver(dir)

	err = saver.Save(signature, manifest, artifacts)
	assert.Nil(t, err)

	loader := release.NewFileSystemLoader(dir)
	sig, mani, arts, err := loader.Load()
	assert.Equal(t, signature, sig)
	assert.Equal(t, manifest, mani)
	//FIXME: improve this shit
	assert.Equal(t, artifacts[0].Digests(), arts[0].Digests())
}
