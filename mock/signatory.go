package mock

import (
	"github.com/awnumar/memguard"
	"github.com/stoic-cli/stoic-release"
)

type signatory struct {
	privateKey *memguard.LockedBuffer
}

func Signatory(privateKey []byte) release.Signatory {
	//FIXME: Will it bite us in the ass not to check this
	//error or will things hopefully start failing?
	k, _ := memguard.NewImmutableFromBytes(privateKey)
	return &signatory{
		privateKey: k,
	}
}

func (s *signatory) PrivateKey() *memguard.LockedBuffer {
	return s.privateKey
}

func ValidSignatory() (release.Signatory, error) {
	pkArm, err := ArmoredToByte(SignerPriv)
	if err != nil {
		return nil, err
	}
	pk, err := memguard.NewImmutableFromBytes(pkArm)
	if err != nil {
		return nil, err
	}
	return release.NewSignatory(pk), nil
}