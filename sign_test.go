package release_test

import (
	"testing"
	"time"

	"github.com/stoic-cli/stoic-release"
	"github.com/stoic-cli/stoic-release/mock"
	"github.com/stoic-cli/stoic-release/pgp"
	"github.com/stretchr/testify/assert"
)

func testTime() time.Time {
	return time.Unix(13455533, 23333)
}

var testSignature = `-----BEGIN PGP SIGNATURE-----

wsFcBAABCAAQBQIAzVCtCRAbjALTQVnSbAAAp3kQABwXELCPeJVo7eUGEWG6f55Z
SiYJZNe44ib6A3/YWGf9iC3T0kuXVVf0NE/Jj0JU1Koyvp+vCeD+Ubk/Q+uG5WJH
+a60aoWfxCxFb/XkfB3cNq+5/9sojHYyUnpbYdsBvh0w2qDwRhDe5dZhE2hMqel3
bFMK20v9vIw1sLaAJL/h0cXcDxImDrj6LQBvFQFeVNDUQSXq68fKkVfEFeBr7ygT
VGup7ijFQPA9yqJwzuWEbusD0Nj9dKF81LtUw11rutrLbABGTRZFAoit6Sq5rJcn
w7WL9ojiuO05ETxBLVuAfzoYeh+oYZjMnJ8OxzSu8EEg7Uwey04GCBZM0VTe8oUM
8khf+g3QIGW7DOgCmrFVcwuONoh83S4F+0fY5MJGqJzegen1JufUD3ymVWR6ssvc
mWqpgimkAomoZ24r7x5zcXLAl0kOLE+PbWhivI47ensDn7MIISWojPn3y1wpnUxg
qijspDrJbmL3xKEWj3RIAT5wA/uZbcoe/sETClrsXSzOE+5f6OHPF+T0NrW2fPI+
QVfDLdfuHhOBNkmDYT/WZZ2xffq5nMkn1Bs3s++gLTtqJ0BktTQpE2JE4WE91xOX
VKpFLFIUhdBOPBkwEtcMF1m/V6hgwKf9xS7kHWMgO1yGHZM5DuKj5MW3x9kIrAUt
akckybeawRp5/eFv12z3
=O0Zl
-----END PGP SIGNATURE-----`

func TestNewSigner(t *testing.T) {
	testCases := []struct {
		name      string
		signatory release.Signatory
		sign      []byte
		expect    interface{}
		expectErr bool
	}{
		{
			name:      "Invalid private key",
			signatory: mock.Signatory([]byte("something not valid")),
			expect:    "failed to sign data: failed to read entity: openpgp: invalid data: tag byte does not have MSB set",
			expectErr: true,
		},
		{
			name: "Should match",
			signatory: mock.Signatory(func() []byte {
				pk, err := mock.ArmoredToByte(mock.SignerPriv)
				assert.Nil(t, err)
				return pk
			}()),
			sign:   mock.Signed,
			expect: testSignature,
		},
	}

	config := pgp.DefaultConfig
	config.Time = testTime
	for _, tc := range testCases {
		got, err := release.NewSigner(config).Sign(tc.signatory, tc.sign)
		if tc.expectErr {
			assert.Error(t, err, tc.name)
			assert.Equal(t, err.Error(), tc.expect, tc.name)
			assert.Nil(t, got, tc.name)
		} else {
			assert.Nil(t, err, tc.name)
			assert.Equal(t, tc.expect, string(got), tc.name)
		}
	}
}
