package release

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ManifestLoader provides the interface for loading
// a manifest
type ManifestLoader interface {
	Read(reader io.Reader) (Manifester, error)
}

// Manifester provides the interface for interacting
// with a manifest
type Manifester interface {
	Name() string
	NormalisedName() string
	Version() string
	Artifacts() []ManifestArtifact
	Serialise() (io.Reader, error)
}

// Manifest  contains the data related to a release
type Manifest struct {
	ReleaseName      string             `yaml:"name"`
	ReleaseVersion   string             `yaml:"version"`
	ReleaseSignee    ManifestSignee     `yaml:"signee"`
	ReleaseArtifacts []ManifestArtifact `yaml:"artifacts"`
}

// ManifestSignee contains the identity that signed
// a release
type ManifestSignee struct {
	User string
	Key  string
	Type SigneeType
}

// ManifestArtifact contains the metadata of a
// release artifact
type ManifestArtifact struct {
	Name    string
	Type    ArtifactType
	Digests map[DigestType]string
}

// NewManifestLoader returns a loader for recreating
// a manifest
func NewManifestLoader() ManifestLoader {
	return &Manifest{}
}

// NewManifest creates a new manifest
func NewManifest(projectName string, version string, signee Signee, artifacts []Artifact) Manifester {
	var manifestArtifacts []ManifestArtifact
	for _, a := range artifacts {
		manifestArtifacts = append(manifestArtifacts, ManifestArtifact{
			Name:    a.NormalisedName(version),
			Type:    a.Type(),
			Digests: a.Digests(),
		})
	}
	return &Manifest{
		ReleaseName:    projectName,
		ReleaseVersion: version,
		ReleaseSignee: ManifestSignee{
			User: signee.User(),
			Key:  signee.Key(),
			Type: signee.Type(),
		},
		ReleaseArtifacts: manifestArtifacts,
	}
}

// Serialise encodes a manifest
func (m *Manifest) Serialise() (io.Reader, error) {
	d, err := yaml.Marshal(&m)
	return strings.NewReader(string(d)), err
}

// Read decodes a manifest
func (m *Manifest) Read(reader io.Reader) (Manifester, error) {
	err := yaml.NewDecoder(reader).Decode(m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read manifest")
	}
	return m, nil
}

// NormalisedName returns a normalised version of the name for the manifest
// of a release
func (m *Manifest) NormalisedName() string {
	return fmt.Sprintf("%s_%s.manifest", strings.ToLower(m.ReleaseName), strings.ToLower(m.Version()))
}

// Version returns the version of the release
func (m *Manifest) Version() string {
	return m.ReleaseVersion
}

// Artifacts returns the artifacts of the release
func (m *Manifest) Artifacts() []ManifestArtifact {
	return m.ReleaseArtifacts
}

// Name returns the name of the release
func (m *Manifest) Name() string {
	return m.ReleaseName
}
