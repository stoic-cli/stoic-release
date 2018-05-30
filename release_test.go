package release_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stoic-cli/stoic-release"
	"github.com/stoic-cli/stoic-release/mock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	p := "MyProject"
	v := release.Version(release.NewProvidedVersion(1, 0, 0))
	d := release.NewDigester(release.DigestTypeSHA256, release.DigestTypeSHA1)

	a1, err := release.NewBinaryArtifact(ioutil.NopCloser(strings.NewReader("some content")), p, release.OperatingSystemTypeDarwin, release.ArchTypeamd64)
	assert.Nil(t, err)
	a2, err := release.NewArtifact(ioutil.NopCloser(strings.NewReader("some content")), p, release.ArtifactTypeReleaseNotes)
	assert.Nil(t, err)
	artifacts := []release.Artifact{a1, a2}

	releaser := release.New(p, v).Add(d, artifacts[0], artifacts[1:]...)
	manifest, artifacts, err := releaser.Create(mock.ValidSignee())
	assert.Nil(t, err)
	var buf bytes.Buffer
	m, err := manifest.Serialise()
	assert.Nil(t, err)
	_, err = io.Copy(&buf, m)

	// Sign
	signatory, err := mock.ValidSignatory()
	assert.Nil(t, err)
	signature, err := releaser.Sign(signatory, buf.Bytes())
	assert.Nil(t, err)

	var buf2 bytes.Buffer
	m, err = manifest.Serialise()
	assert.Nil(t, err)
	_, err = io.Copy(&buf2, m)

	identities, err := releaser.VerifySignature(mock.ValidSignee(), buf2.Bytes(), signature)
	assert.Nil(t, err)
	fmt.Println(identities)
}
