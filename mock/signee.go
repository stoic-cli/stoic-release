package mock

import "github.com/stoic-cli/stoic-release"

type signee struct {
	publicKey  []byte
	signeeType release.SigneeType
	key        string
	user       string
	err        error
}

func ValidSignee() release.Signee {
	return Signee(SignerKeyID, "bob", release.GithubSigneeType, []byte(SignerPub), nil)
}

func Signee(key, user string, signeeType release.SigneeType, publicKey []byte, err error) release.Signee {
	return &signee{
		key:        key,
		user:       user,
		signeeType: signeeType,
		publicKey:  publicKey,
		err:        err,
	}
}

func (s *signee) PublicKey() ([]byte, error) {
	return s.publicKey, s.err
}

func (s *signee) Key() string {
	return s.key
}

func (s *signee) User() string {
	return s.user
}

func (s *signee) Type() release.SigneeType {
	return s.signeeType
}
