package release

import (
	"github.com/pkg/errors"
	"github.com/stoic-cli/stoic-release/pgp"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"fmt"
)

// Verifier provides the interface required for verifying
// a release
type Verifier interface {
	// VerifySignature asserts that the public key of the signee is available
	// and that the key verifies the signature
	VerifySignature(signee Signee, signed []byte, signature []byte) ([]string, error)

	// VerifyDigests asserts that the provided hashes match that of the given
	// artifact
	VerifyDigests(digests map[DigestType]string, reader io.Reader) error
}

// NewVerifier creates a new stand alone verifier
func NewVerifier(config *packet.Config) Verifier {
	return &verifier{
		config: config,
	}
}

type verifier struct {
	config *packet.Config
}

func (v *verifier) VerifySignature(signee Signee, signed []byte, signature []byte) ([]string, error) {
	signeeKey, err := signee.PublicKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch signee's public key")
	}
	if v.config == nil {
		v.config = pgp.DefaultConfig
	}
	identities, err := pgp.Verify(signeeKey, signed, signature, v.config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify signature")
	}
	return identities, nil
}

var (
	ErrNoDigests = errors.New("no digests provided")
)

func (v *verifier) VerifyDigests(digests map[DigestType]string, reader io.Reader) error {
	if len(digests) == 0 {
		return ErrNoDigests
	}

	var digesters []DigestType
	for d := range digests {
		digesters = append(digesters, d)
	}

	var digester Digester
	if len(digesters) == 1 {
		digester = NewDigester(digesters[0])
	} else {
		digester = NewDigester(digesters[0], digesters[1:]...)
	}

	ourDigests, err := digester.Digest(reader)
	if err != nil {
		return errors.Wrap(err, "failed to verify digests")
	}

	for hashType, digest := range ourDigests {
		if digests[hashType] != digest {
			return fmt.Errorf("verification failed, hash mismatch, got: %s, expected: %s", digest, digests[hashType])
		}
		delete(ourDigests, hashType)
	}

	if len(ourDigests) != 0 { // FIXME: meh, should never reach this
		return fmt.Errorf("failed to verify all hashes we produced")
	}

	return nil
}
