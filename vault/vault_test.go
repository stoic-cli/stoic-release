package vault

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/awnumar/memguard"
)

func LockedBuffer(t *testing.T, buf string) *memguard.LockedBuffer {
	lbuf, err := memguard.NewImmutableFromBytes([]byte(buf))
	assert.Nil(t, err)
	return lbuf
}

func TestVault(t *testing.T) {
	testCases := []struct {
		name      string
		sealPass  string
		openPass  string
		message   string
		expectErr bool
	}{
		{
			name:     "Correct pass",
			sealPass: "secret1",
			openPass: "secret1",
			message:  "this is my message",
		},
		{
			name:      "Incorrect pass",
			sealPass:  "secret1",
			openPass:  "breaking",
			message:   "this is my other message",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		sealed, err := Seal(LockedBuffer(t, tc.sealPass), LockedBuffer(t, tc.message))
		assert.Nil(t, err, tc.name)
		got, err := Open(LockedBuffer(t, tc.openPass), sealed)
		if tc.expectErr {
			assert.Error(t, err, tc.name)
		} else {
			assert.Equal(t, string(got.Buffer()), tc.message, tc.name)
		}
	}
}
