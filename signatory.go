package release

import "github.com/awnumar/memguard"

// Signatory provides the interface required
// to sign a release
type Signatory interface {
	PrivateKey() *memguard.LockedBuffer
}

type signatory struct {
	privateKey *memguard.LockedBuffer
}

// NewSignatory creates a new Signatory
func NewSignatory(privateKey *memguard.LockedBuffer) Signatory {
	return &signatory{
		privateKey: privateKey,
	}
}

// PrivateKey returns the private key of the signatory
func (s *signatory) PrivateKey() *memguard.LockedBuffer {
	return s.privateKey
}
