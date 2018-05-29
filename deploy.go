package release

type Deployer interface {
	// Deploy the artifacts to the given destination
	Deploy(signature []byte, manifest Manifester, artifacts []Artifact) error
}

type deployer struct {
	deployers []Deployer
}

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

func NewGithubDeployer(user, token string) Deployer {
	return &githubDeployer{
		user: user,
		token: token,
	}
}

func (gd *githubDeployer) Deploy(signature []byte, manifester Manifester, artifacts []Artifact) error {
	return nil
}
