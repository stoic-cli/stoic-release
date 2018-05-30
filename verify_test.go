package release_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stoic-cli/stoic-release"
	"github.com/stoic-cli/stoic-release/mock"
	"github.com/stoic-cli/stoic-release/pgp"
	"github.com/stretchr/testify/assert"
)

func TestVerifySignature(t *testing.T) {
	testCases := []struct {
		name      string
		signee    release.Signee
		signed    []byte
		signature []byte
		expect    interface{}
		expectErr bool
	}{
		{
			name:      "Signee with err",
			signee:    mock.Signee("", "", "", nil, fmt.Errorf("some error")),
			expect:    "failed to fetch signee's public key: some error",
			expectErr: true,
		},
		{
			name:      "Signee with wrong public key",
			signee:    mock.Signee("", "", "", mock.AltSignerPub, nil),
			signed:    mock.Signed,
			signature: mock.Signature,
			expect:    "failed to verify signature: failed to check armored detached signature: openpgp: signature made by unknown entity",
			expectErr: true,
		},
		{
			name:      "Signee with correct public key",
			signee:    mock.Signee("", "", "", mock.SignerPub, nil),
			signed:    mock.Signed,
			signature: mock.Signature,
			expect:    []string{"Bob the Builder (I build stuff) <bob@builder.com>"},
		},
	}

	for _, tc := range testCases {
		got, err := release.NewVerifier(pgp.DefaultConfig).VerifySignature(tc.signee, tc.signed, tc.signature)
		if tc.expectErr {
			assert.Error(t, err, tc.name)
			assert.Equal(t, err.Error(), tc.expect, tc.name)
			assert.Nil(t, got, tc.name)
		} else {
			assert.Nil(t, err, tc.name)
			assert.Equal(t, tc.expect, got, tc.name)
		}
	}
}

func TestVerifyDigests(t *testing.T) {
	testCases := []struct {
		name      string
		artifact  io.Reader
		digests   map[release.DigestType]string
		expect    string
		expectErr bool
	}{
		{
			name: "Should work",
			digests: map[release.DigestType]string{
				release.DigestTypeMD5:    "736db904ad222bf88ee6b8d103fceb8e",
				release.DigestTypeSHA1:   "5ec1a3cb71c75c52cf23934b137985bd2499bd85",
				release.DigestTypeSHA256: "373993310775a34f5ad48aae265dac65c7abf420dfbaef62819e2cf5aafc64ca",
				release.DigestTypeSHA512: "47bb28d146567b3be18d06d8468aaa8222183fe6b2a942b17b6a48bbc32bda7213f7dc1acf36677f7710cffa7add3f3656597630bf0d591f34145015f59724e1",
			},
			artifact: strings.NewReader("this is some content"),
		},
		{
			name: "Wrong digest",
			digests: map[release.DigestType]string{
				release.DigestTypeMD5: "736db904ad222b",
			},
			artifact:  strings.NewReader("this is some content"),
			expect:    "verification failed, hash mismatch, got: 736db904ad222bf88ee6b8d103fceb8e, expected: 736db904ad222b",
			expectErr: true,
		},
		{
			name:      "No digests",
			digests:   nil,
			expect:    release.ErrNoDigests.Error(),
			expectErr: true,
		},
		{
			name: "Nil reader",
			digests: map[release.DigestType]string{
				release.DigestTypeSHA1: "5ec1a3cb71c75c52cf23934b137985bd2499bd85",
			},
			artifact:  nil,
			expect:    "failed to verify digests: reader is nil",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		err := release.NewVerifier(pgp.DefaultConfig).VerifyDigests(tc.digests, tc.artifact)
		if tc.expectErr {
			assert.Equal(t, tc.expect, err.Error(), tc.name)
		} else {
			assert.Nil(t, err)
		}
	}
}
