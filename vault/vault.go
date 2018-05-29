package vault

import (
	"errors"
	"github.com/awnumar/memguard"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/scrypt"
	"unsafe"
)

// Based on the code from: https://leanpub.com/gocrypto/read
// and modified to use mempage.LockedBuffer

const (
	keySize    = 32
	nonceSize  = 24
	saltSize   = 32
	iterations = 1048576
)

var (
	ErrEncrypt = errors.New("encryption failed")
	ErrDecrypt = errors.New("decryption failed")
	ErrCopy    = errors.New("copy failed")
)

func encrypt(key *memguard.LockedBuffer, message *memguard.LockedBuffer) ([]byte, error) {
	nonce, err := memguard.NewImmutableRandom(nonceSize)
	if nonce != nil {
		defer nonce.Destroy()
	}
	if err != nil {
		return nil, ErrEncrypt
	}

	out := make([]byte, nonce.Size())
	if n := copy(out, nonce.Buffer()); n != nonceSize {
		return nil, ErrEncrypt
	}

	nonceArrayPtr := (*[nonceSize]byte)(unsafe.Pointer(&nonce.Buffer()[0]))
	keyArrayPtr := (*[keySize]byte)(unsafe.Pointer(&key.Buffer()[0]))

	out = secretbox.Seal(out, message.Buffer(), nonceArrayPtr, keyArrayPtr)
	return out, nil
}

func decrypt(key *memguard.LockedBuffer, message []byte) (*memguard.LockedBuffer, error) {
	if len(message) < (nonceSize + secretbox.Overhead) {
		return nil, ErrDecrypt
	}

	var nonce [nonceSize]byte
	copy(nonce[:], message[:nonceSize])

	keyArrayPtr := (*[keySize]byte)(unsafe.Pointer(&key.Buffer()[0]))
	out, ok := secretbox.Open(nil, message[nonceSize:], &nonce, keyArrayPtr)
	if !ok {
		return nil, ErrDecrypt
	}

	guarded, err := memguard.NewImmutableFromBytes(out)
	if err != nil {
		memguard.WipeBytes(out)
		if guarded != nil {
			guarded.Destroy()
		}
		return nil, ErrDecrypt
	}

	return guarded, nil
}

func deriveKey(pass, salt *memguard.LockedBuffer) (*memguard.LockedBuffer, error) {
	// scrypt will segfault due to the memory
	// guards provided by memguard unless we
	// create a copy here
	passBuff := make([]byte, pass.Size())
	defer memguard.WipeBytes(passBuff)
	if copy(passBuff, pass.Buffer()) != pass.Size() {
		return nil, ErrCopy
	}

	key, _ := scrypt.Key(passBuff, salt.Buffer(), iterations, 8, 1, keySize)
	memguard.WipeBytes(passBuff)

	return memguard.NewImmutableFromBytes(key) // This also wipes the key slice
}

// Seal encrypts the provided message using NaCl, we accept a LockedBuffer
// because we expect the object to be sealed should be kept as safe as possible
func Seal(pass *memguard.LockedBuffer, message *memguard.LockedBuffer) ([]byte, error) {
	salt, err := memguard.NewImmutableRandom(saltSize)
	if salt != nil {
		defer salt.Destroy()
	}
	if err != nil {
		return nil, ErrEncrypt
	}

	key, err := deriveKey(pass, salt)
	if key != nil {
		defer key.Destroy()
	}
	if err != nil {
		return nil, ErrEncrypt
	}
	out, err := encrypt(key, message)
	key.Destroy()
	if err != nil {
		return nil, ErrEncrypt
	}

	out = append(salt.Buffer(), out...)
	return out, nil
}

const overhead = saltSize + secretbox.Overhead + nonceSize

// Open decrypts the provided message using NaCl, we return a LockedBuffer
// because we expect anything that is stored in an encrypted fashion
// should be kept as secure as possible
func Open(pass *memguard.LockedBuffer, message []byte) (*memguard.LockedBuffer, error) {
	if len(message) < overhead {
		return nil, ErrDecrypt
	}

	salt, err := memguard.NewImmutableFromBytes(message[:saltSize])
	if salt != nil {
		defer salt.Destroy()
	}
	if err != nil {
		return nil, ErrDecrypt
	}

	key, err := deriveKey(pass, salt)
	if key != nil {
		defer key.Destroy()
	}
	if err != nil {
		return nil, ErrDecrypt
	}

	out, err := decrypt(key, message[saltSize:])
	key.Destroy()
	if err != nil {
		return nil, ErrDecrypt
	}

	return out, nil
}
