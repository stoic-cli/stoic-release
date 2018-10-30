package release

// Deployer defines the operations required for deploying
// a release
type Deployer interface {
	// Deploy the artifacts to the given destination
	Deploy(signature []byte, manifest Manifester, artifacts []Artifact) error
}

type deployer struct {
	deployers []Deployer
}

// NewDeployers creates a stand-alone deployer that
// composes multiple deployers
func NewDeployers(deployers []Deployer) Deployer {
	return &deployer{
		deployers: deployers,
	}
}

// Deploy the manifest, artifacts and signature
func (d *deployer) Deploy(signature []byte, manifester Manifester, artifacts []Artifact) error {
	for _, deps := range d.deployers {
		err := deps.Deploy(signature, manifester, artifacts)
		if err != nil {
			return err
		}
	}
	return nil
}

type githubDeployer struct {
	user  string
	token string
}

// NewGithubDeployer creates a deployer that can
// push to github
func NewGithubDeployer(user, token string) Deployer {
	return &githubDeployer{
		user:  user,
		token: token,
	}
}

// Deploy the provided data
func (gd *githubDeployer) Deploy(signature []byte, manifester Manifester, artifacts []Artifact) error {
	return nil
}
