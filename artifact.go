package release

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"strings"
)

// Artifact provides the interface for
// interacting with a given artifact
type Artifact interface {
	NormalisedName(version string) string
	Type() ArtifactType
	Digests() map[DigestType]string
	SetDigests(digests map[DigestType]string)
	Content() io.Reader
}

// ArtifactType enumerates the available artifact
// types
type ArtifactType string

const (
	ArtifactTypeReleaseNotes = "relnotes"
	ArtifactTypeChangeSet    = "changeset"
	ArtifactTypeReadme       = "readme"
	ArtifactTypeBinary       = "bin"
	ArtifactTypeDEB          = "deb"
	ArtifactTypeRPM          = "rpm"
)

type normaliseNameFn func(version string) string

type artifact struct {
	normaliseNameFn normaliseNameFn
	artifactType    ArtifactType
	digests         map[DigestType]string
	content         []byte
}

func readContent(content io.ReadCloser) ([]byte, error) {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, content)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read content")
	}
	err = content.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close content")
	}
	return buf.Bytes(), nil
}

func newArtifact(content io.ReadCloser, name string, artifactType ArtifactType, fn normaliseNameFn) (Artifact, error) {
	c, err := readContent(content)
	if err != nil {
		return nil, err
	}

	return &artifact{
		normaliseNameFn: fn,
		artifactType:    artifactType,
		digests:         map[DigestType]string{},
		content:         c,
	}, nil
}

func NewArtifact(content io.ReadCloser, projectName string, artifactType ArtifactType) (Artifact, error) {
	return newArtifact(content, projectName, artifactType, func(version string) string {
		return fmt.Sprintf("%s_%s.%s", strings.ToLower(projectName), strings.ToLower(version), artifactType)
	})
}

func NewBinaryArtifact(content io.ReadCloser, projectName string, os OperatingSystemType, arch ArchType) (Artifact, error) {
	return newArtifact(content, projectName, ArtifactTypeBinary, func(version string) string {
		return fmt.Sprintf("%s_%s-%s.%s.%s", strings.ToLower(projectName), strings.ToLower(version), os, arch, ArtifactTypeBinary)
	})
}

func (a *artifact) NormalisedName(version string) string {
	return a.normaliseNameFn(version)
}

func (a *artifact) Type() ArtifactType {
	return a.artifactType
}

func (a *artifact) SetDigests(digests map[DigestType]string) {
	a.digests = digests
}

func (a *artifact) Digests() map[DigestType]string {
	return a.digests
}

func (a *artifact) Content() io.Reader {
	return bytes.NewReader(a.content)
}
