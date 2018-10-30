package release

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameNormalising(t *testing.T) {
	testCases := []struct {
		name     string
		artifact Artifact
		version  string
		expect   string
	}{
		{
			name: "binary artifact",
			artifact: func() Artifact {
				a, err := NewBinaryArtifact(ioutil.NopCloser(strings.NewReader("some content")), "MyProject", OperatingSystemTypeDarwin, ArchTypeamd64)
				assert.Nil(t, err)
				return a
			}(),
			version: "v1.0.0",
			expect:  "myproject_v1.0.0-darwin.amd64.bin",
		},
		{
			name: "other artifact",
			artifact: func() Artifact {
				a, err := NewArtifact(ioutil.NopCloser(strings.NewReader("some content")), "MyProject", ArtifactTypeReleaseNotes)
				assert.Nil(t, err)
				return a
			}(),
			version: "v1.2.0",
			expect:  "myproject_v1.2.0.relnotes",
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.artifact.NormalisedName(tc.version), tc.expect)
	}
}
