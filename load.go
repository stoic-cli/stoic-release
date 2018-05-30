package release

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"github.com/stoic-cli/stoic-release/pgp"
)

// Loader defines the operations required for
// loading a created release
type Loader interface {
	Load() (signature []byte, manifester Manifester, artifacts []Artifact, err error)
}

// LoadFinaliser defines the operations required
// loading and finalising deployment of a release
type LoadFinaliser interface {
	Loader
	Finaliser
}

type loadFinaliser struct {
	deployers []Deployer
	Loader
	Verifier
	Deployer
}

// NewLoader creates a new loader
func NewLoader(loader Loader, deployer Deployer, deployers ...Deployer) (LoadFinaliser, error) {
	return &loadFinaliser{
		Loader:   loader,
		Verifier: NewVerifier(pgp.DefaultConfig),
		Deployer: NewDeployers(append(deployers, deployer)),
	}, nil
}

type fileSystemLoader struct {
	directory string
}

// NewFileSystemLoader creates a loader that can read
// a release from a filesystem
func NewFileSystemLoader(directory string) Loader {
	return &fileSystemLoader{
		directory: directory,
	}
}

func (fs *fileSystemLoader) Load() ([]byte, Manifester, []Artifact, error) {
	absPath, err := filepath.Abs(fs.directory)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to get absolute path")
	}

	manifestFile, err := fileFromGlob(absPath, "*.manifest")
	if err != nil {
		return nil, nil, nil, err
	}
	manifester, err := NewManifestLoader().Read(manifestFile)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to load manifest")
	}

	var signature bytes.Buffer
	manifestSigFile, err := fileFromGlob(absPath, "*.manifest.asc")
	if err != nil {
		return nil, nil, nil, err
	}
	_, err = io.Copy(&signature, manifestSigFile)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to load manifest signature")
	}

	var artifacts []Artifact
	manifestArtifacts := manifester.Artifacts()
	for _, artifact := range manifestArtifacts {
		f, err := os.Open(path.Join(absPath, artifact.Name))
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "failed to load artifact: %s", artifact.Name)
		}

		var art Artifact
		//FIXME: Don't do this here..
		//FIXME: Don't do it this way.. serialise per artifact type..
		switch artifact.Type {
		case ArtifactTypeBinary:
			binRegex := regexp.MustCompile(".*-(?P<os>[a-z]+).(?P<arch>[a-z0-9]+).bin$")
			result := make(map[string]string)
			match := binRegex.FindStringSubmatch(artifact.Name)
			for i, name := range binRegex.SubexpNames() {
				if i != 0 && name != "" {
					result[name] = match[i]
				}
			}
			art, err = NewBinaryArtifact(f, manifester.Name(), OperatingSystemType(result["os"]), ArchType(result["arch"]))
		default:
			art, err = NewArtifact(f, manifester.Name(), artifact.Type)
		}
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "failed to recreate artifact")
		}
		art.SetDigests(artifact.Digests)
		artifacts = append(artifacts, art)
	}

	return signature.Bytes(), manifester, artifacts, nil
}

func fileFromGlob(basePath, pattern string) (*os.File, error) {
	matches, err := filepath.Glob(path.Join(basePath, pattern))
	if err != nil {
		return nil, errors.Wrap(err, "failed to glob filesystem")
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("found too many matches for: %s, expected: %d, got: %d", pattern, 1, len(matches))
	}
	file, err := os.Open(matches[0])
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", matches[0])
	}
	return file, nil
}
