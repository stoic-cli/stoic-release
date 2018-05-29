package pgp_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stoic-cli/stoic-release/mock"
	"github.com/stoic-cli/stoic-release/pgp"
)

func TestHappyPath(t *testing.T) {
	// Signer
	kp, err := pgp.NewSigner("Bob the Builder", "I build stuff", "bob@builder.com", pgp.DefaultConfig)
	assert.Nil(t, err)

	// Sign
	msg := []byte("This is my message\n")
	signed, err := pgp.Sign(kp.PrivateKey, msg, pgp.DefaultConfig)
	assert.Nil(t, err)

	// Verify
	identities, err := pgp.Verify(kp.PublicKey, msg, signed, pgp.DefaultConfig)
	assert.Nil(t, err)
	assert.Len(t, identities, 1)
	assert.Equal(t, identities[0], "Bob the Builder (I build stuff) <bob@builder.com>")
}

func TestVerifySignature(t *testing.T) {
	_, err := pgp.Verify(mock.SignerPub, mock.Signed, mock.Signature, pgp.DefaultConfig)
	assert.Nil(t, err)

	_, err = pgp.Verify(mock.AltSignerPub, mock.Signed, mock.Signature, pgp.DefaultConfig)
	assert.Error(t, err)
}

//func TestMeh(t *testing.T) {
//	pk, _ := mock.ArmoredToByte(mock.SignerPriv)
//	entity, _ := openpgp.ReadEntity(packet.NewReader(bytes.NewReader(pk)))
//	fmt.Printf("%x\n", entity.PrivateKey.Fingerprint)
//	fmt.Printf("%x", entity.PrivateKey.KeyId)
//}