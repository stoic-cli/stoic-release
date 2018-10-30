package release

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

// Saver provides an interface for storing a release
type Saver interface {
	Save(signature []byte, manifest Manifester, artifacts []Artifact) error
}

type saver struct {
	savers []Saver
}

// NewSavers returns a composition for multiple
// savers
func NewSavers(savers []Saver) Saver {
	return &saver{
		savers: savers,
	}
}

// Save the manifest, artifacts and signature
func (s *saver) Save(signature []byte, manifest Manifester, artifacts []Artifact) error {
	for _, saver := range s.savers {
		err := saver.Save(signature, manifest, artifacts)
		if err != nil {
			return errors.Wrap(err, "failed to save")
		}
	}
	return nil
}

type fileSystemSaver struct {
	directory string
}

// NewFileSystemSaver will store the provided release as files
// in a directory
func NewFileSystemSaver(directory string) Saver {
	return &fileSystemSaver{
		directory: directory,
	}
}

// Save the release to the file system
func (fs *fileSystemSaver) Save(signature []byte, manifest Manifester, artifacts []Artifact) error {
	absPath, err := filepath.Abs(fs.directory)
	if err != nil {
		return errors.Wrap(err, "failed to get absolute path")
	}

	err = os.MkdirAll(absPath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	serialisedManifest, err := manifest.Serialise()
	if err != nil {
		return errors.Wrap(err, "failed to serialise manifest")
	}
	err = createAndWriteFile(serialisedManifest, absPath, manifest.NormalisedName())
	if err != nil {
		return err
	}

	err = createAndWriteFile(bytes.NewReader(signature), absPath, fmt.Sprintf("%s.asc", manifest.NormalisedName()))
	if err != nil {
		return err
	}

	v := manifest.Version()
	for _, artifact := range artifacts {
		err = createAndWriteFile(artifact.Content(), absPath, artifact.NormalisedName(v))
		if err != nil {
			return err
		}
	}

	return nil
}

func createAndWriteFile(content io.Reader, basePath, name string) error {
	fileName := path.Join(basePath, name)
	file, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create file: %s", name))
	}
	_, err = io.Copy(file, content)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write content to file: %s", name))
	}
	err = file.Close()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to close file: %s", name))
	}
	return nil
}
