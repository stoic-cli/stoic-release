package release

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
)

type ManifestLoader interface {
	Read(reader io.Reader) (Manifester, error)
}

type Manifester interface {
	Name() string
	NormalisedName() string
	Version() string
	Artifacts() []ManifestArtifact
	Serialise() (io.Reader, error)
}

type Manifest struct {
	ReleaseName      string             `yaml:"name"`
	ReleaseVersion   string             `yaml:"version"`
	ReleaseSignee    ManifestSignee     `yaml:"signee"`
	ReleaseArtifacts []ManifestArtifact `yaml:"artifacts"`
}

type ManifestSignee struct {
	User string
	Key  string
	Type SigneeType
}

type ManifestArtifact struct {
	Name    string
	Type    ArtifactType
	Digests map[DigestType]string
}

func NewManifestLoader() ManifestLoader {
	return &Manifest{}
}

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

func (m *Manifest) Serialise() (io.Reader, error) {
	d, err := yaml.Marshal(&m)
	return strings.NewReader(string(d)), err
}

func (m *Manifest) Read(reader io.Reader) (Manifester, error) {
	err := yaml.NewDecoder(reader).Decode(m)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read manifest")
	}
	return m, nil
}

func (m *Manifest) NormalisedName() string {
	return fmt.Sprintf("%s_%s.manifest", strings.ToLower(m.ReleaseName), strings.ToLower(m.Version()))
}

func (m *Manifest) Version() string {
	return m.ReleaseVersion
}

func (m *Manifest) Artifacts() []ManifestArtifact {
	return m.ReleaseArtifacts
}

func (m *Manifest) Name() string {
	return m.ReleaseName
}
