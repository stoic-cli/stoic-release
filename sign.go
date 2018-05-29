package release

import (
	"github.com/stoic-cli/stoic-release/pgp"
	"github.com/pkg/errors"
	"golang.org/x/crypto/openpgp/packet"
)

// Signer provides the interface required to sign data
type Signer interface {
	// Sign uses the signatory to sign the provided data
	Sign(signatory Signatory, sign []byte) (signed []byte, err error)
}

type signer struct {
	config *packet.Config
}

//FIXME: Don't like this, required?
func NewSigner(config *packet.Config) Signer {
	return &signer{
		config: config,
	}
}

// Sign the provided artifact
func (s *signer) Sign(signatory Signatory, sign []byte) ([]byte, error) {
	if s.config == nil {
		s.config = pgp.DefaultConfig
	}
	signed, err := pgp.Sign(signatory.PrivateKey(), sign, s.config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign data")
	}
	return signed, nil
}
