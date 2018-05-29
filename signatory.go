package release

import "github.com/awnumar/memguard"

type Signatory interface {
	PrivateKey() *memguard.LockedBuffer
}

type signatory struct {
	privateKey *memguard.LockedBuffer
}

func NewSignatory(privateKey *memguard.LockedBuffer) Signatory {
	return &signatory{
		privateKey: privateKey,
	}
}

func (s *signatory) PrivateKey() *memguard.LockedBuffer {
	return s.privateKey
}
