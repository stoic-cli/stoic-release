package pgp

import (
	"bytes"
	"crypto"
	"time"

	"github.com/awnumar/memguard"
	"github.com/pkg/errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

// DefaultBits is the default size used when
// creating a new key
var DefaultBits = 4096

// DefaultCipher sets the default cipher algorithm
var DefaultCipher = packet.CipherAES256

// DefaultHash sets the default hash algorithm
var DefaultHash = crypto.SHA256

// DefaultTime sets the default time function
var DefaultTime = func() time.Time {
	return time.Now()
}

// DefaultConfig is the default config used
// creating a new key
var DefaultConfig = &packet.Config{
	Time:          DefaultTime,
	DefaultHash:   DefaultHash,
	DefaultCipher: DefaultCipher,
	RSABits:       DefaultBits,
}

// KeyPair contains a pgp keypair that can be used
// to make and sign artifacts
type KeyPair struct {
	PublicKey            []byte
	PublicKeyFingerPrint [20]byte
	PublicKeyID          uint64
	PrivateKey           *memguard.LockedBuffer
}

func encodePrivateKey(e *openpgp.Entity, cfg *packet.Config) (*memguard.LockedBuffer, error) {
	var privKey bytes.Buffer
	var privKeyGuarded *memguard.LockedBuffer

	err := e.SerializePrivate(&privKey, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialise private key")
	}

	privKeyGuarded, err = memguard.NewImmutableFromBytes(privKey.Bytes())
	if err != nil {
		if privKeyGuarded != nil {
			privKeyGuarded.Destroy()
		}
		return nil, errors.Wrap(err, "failed to protect private key")
	}

	return privKeyGuarded, nil
}

func encodePublicKey(e *openpgp.Entity) ([]byte, error) {
	var pubKey, pubKeyArmored bytes.Buffer
	err := e.Serialize(&pubKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialise public key")
	}

	w, err := armor.Encode(&pubKeyArmored, openpgp.PublicKeyType, map[string]string{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to armor public key")
	}
	_, err = w.Write(pubKey.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "failed to write armored public key")
	}
	err = w.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close armored public key")
	}

	return pubKeyArmored.Bytes(), nil
}

// NewSigner creates a new pgp keypair capable of signing artifacts
// using the provided input parameters
func NewSigner(name, comment, email string, config *packet.Config) (*KeyPair, error) {
	entity, err := openpgp.NewEntity(name, comment, email, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new entity")
	}

	privKeyGuarded, err := encodePrivateKey(entity, DefaultConfig)
	if err != nil {
		return nil, err
	}

	pubKeyArmored, err := encodePublicKey(entity)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		PublicKey:            pubKeyArmored,
		PublicKeyFingerPrint: entity.PrimaryKey.Fingerprint,
		PublicKeyID:          entity.PrimaryKey.KeyId,
		PrivateKey:           privKeyGuarded,
	}, nil
}

// Sign loads a signing entity from the provided signer and uses it
// to create an armored detached signature
func Sign(signer *memguard.LockedBuffer, sign []byte, config *packet.Config) ([]byte, error) {
	entity, err := openpgp.ReadEntity(packet.NewReader(bytes.NewReader(signer.Buffer())))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read entity")
	}

	var signed bytes.Buffer
	err = openpgp.ArmoredDetachSign(&signed, entity, bytes.NewReader(sign), config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign message")
	}

	return signed.Bytes(), nil
}

// Verify loads a signers public key and verify the signed object using
// the provided signature
// FIXME: We probably want to structure this part to be more flexible, e.g.,
// by making it possible to read unarmored key rings and signatures, etc.
func Verify(signer []byte, signed []byte, signature []byte, config *packet.Config) ([]string, error) {
	keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(signer))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read armored keyring")
	}

	entity, err := openpgp.CheckArmoredDetachedSignature(keyring, bytes.NewReader(signed), bytes.NewReader(signature))
	if err != nil {
		return nil, errors.Wrap(err, "failed to check armored detached signature")
	}

	var identities []string
	for _, identity := range entity.Identities {
		identities = append(identities, identity.Name)
	}

	return identities, nil
}
