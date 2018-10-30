package release

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
)

func TestNewProvidedVersion(t *testing.T) {
	expect := "v1.0.0"
	v := NewProvidedVersion(1, 0, 0)
	got, err := v.Version()
	assert.Equal(t, expect, got)
	assert.Nil(t, err)
}

func TestNewGitVersion(t *testing.T) {
	testCases := []struct {
		name        string
		url         string
		branch      string
		major       int
		expect      interface{}
		expectError bool
	}{
		{
			name:   "Simple",
			url:    "https://github.com/src-d/go-siva",
			branch: "master",
			major:  1,
			expect: "v1.17.53",
		},
		{
			name:   "Different branch",
			url:    "https://github.com/paulbes/mergedcallbacks",
			branch: "master",
			major:  1,
			expect: "v1.0.3",
		},
	}

	for _, tc := range testCases {
		dir, err := ioutil.TempDir("", "")
		assert.Nil(t, err)
		_, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL: tc.url,
		})
		assert.Nil(t, err)
		v := NewGitHistoryVersion(dir, tc.branch, tc.major)
		got, err := v.Version()
		if tc.expectError {
			assert.Equal(t, got, "")
			assert.Equal(t, err, tc.expect)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, got, tc.expect)
		}
	}
}
