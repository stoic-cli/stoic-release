package release_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"strings"
	"github.com/stoic-cli/stoic-release/mock"
	"github.com/stoic-cli/stoic-release"
)

func TestNewSignee(t *testing.T) {
	testCases := []struct {
		name       string
		user       string
		key        string
		signeeType release.SigneeType
		expect     string
		expectErr  bool
	}{
		//FIXME: Do this with mocks instead of towards actual server..
		{
			name:       "Github with valid user and key",
			user:       mock.GithubPublicKeyUser,
			key:        mock.GithubPublicKeyID,
			signeeType: release.GithubSigneeType,
			expect:     mock.GithubPublicKey,
		},
	}

	for _, tc := range testCases {
		s := release.NewSignee(tc.user, tc.key, tc.signeeType)
		got, err := s.PublicKey()
		if tc.expectErr {
			assert.Error(t, err, tc.name)
			assert.Nil(t, got, tc.name)
		} else {
			assert.Equal(t, tc.expect, string(strings.Replace(string(got), "\r", "", -1)), tc.name)
			assert.Nil(t, err, tc.name)
		}
	}
}

