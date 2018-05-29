package release

import (
	"github.com/pkg/errors"
	"github.com/stoic-cli/stoic-release/pgp"
)

// Releaser provides an interface for all activities
// related to creating, verifying and deploying a
// release
type Releaser interface {
	Adder
	Creator
	Signer
	Saver
	Finaliser
}

// Creator provides the interface required for creating
// a new release
type Creator interface {
	// Create a new release by normalising all artifacts,
	// generating digests and creating a manifest. This
	// will return the manifest and its artifacts if successful
	// otherwise an error.
	Create(signee Signee) (manifest Manifester, artifacts []Artifact, err error)
}

// Adder provides an interface for adding artifacts
// to a release. It also requires a digester
// to be provided.
type Adder interface {
	Add(digester Digester, artifact Artifact, artifacts ...Artifact) Releaser
}

// Finaliser provides an interface for
type Finaliser interface {
	Verifier
	Deployer
}

// New creates a new release
func New(name string, options ...Option) Releaser {
	// Add some sensible default if otherwise not provided
	r := &releaser{
		name:    name,
		version: NewGitHistoryVersion(".", "master", 1),
		Signer: NewSigner(pgp.DefaultConfig),
		Verifier: NewVerifier(pgp.DefaultConfig),
	}
	for _, o := range options {
		o(r)
	}

	r.Deployer = NewDeployers(r.deployers)
	r.Saver = NewSavers(r.savers)

	return r
}

// Deploy adds a deployer
func Deploy(deployer Deployer) Option {
	return func(args *releaser) {
		args.deployers = append(args.deployers, deployer)
	}
}

// Save adds a saver
func Save(saver Saver) Option {
	return func(args *releaser) {
		args.savers = append(args.savers, saver)
	}
}

// Version adds a versioner
func Version(versioner Versioner) Option {
	return func(args *releaser) {
		args.version = versioner
	}
}

type releaser struct {
	name      string
	version   Versioner
	artifacts []Artifact
	digester  Digester
	savers    []Saver
	deployers []Deployer

	// Pull in some external functionality
	Saver
	Deployer
	Verifier
	Signer
}

// Option is the interface required
// for adding an item to the releaser
type Option func(*releaser)

// Add a digester and artifacts to the release
func (o *releaser) Add(digester Digester, artifact Artifact, artifacts ...Artifact) Releaser {
	o.digester = digester
	o.artifacts = append(artifacts, artifact)
	return o
}

// Create a manifest of the release artifacts, including adding
// information on the signing party and digests of the artifacts
func (o *releaser) Create(signee Signee) (Manifester, []Artifact, error) {
	for _, artifact := range o.artifacts {
		digests, err := o.digester.Digest(artifact.Content())
		if err != nil {
			return nil, nil, errors.Wrap(err, "create failed")
		}
		artifact.SetDigests(digests)
	}

	version, err := o.version.Version()
	if err != nil {
		return nil, nil, errors.Wrap(err, "create failed")
	}
	manifest := NewManifest(o.name, version, signee, o.artifacts)

	return manifest, o.artifacts, nil
}
