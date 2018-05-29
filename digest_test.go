package release

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestDigest(t *testing.T) {
	testCases := []struct {
		name      string
		content   io.Reader
		digester  Digester
		expect    interface{}
		expectErr bool
	}{
		{
			name:      "Digester that doesn't exist",
			digester:  NewDigester(DigestType("something wrong")),
			content:   strings.NewReader(""),
			expect:    fmt.Errorf("unsupported digester: something wrong"),
			expectErr: true,
		},
		{
			name:      "Nil reader",
			digester:  NewDigester(DigestTypeSHA512),
			content:   nil,
			expect:    fmt.Errorf("reader is nil"),
			expectErr: true,
		},
		{
			name:      "Empty content",
			digester:  NewDigester(DigestTypeMD5),
			content:   strings.NewReader(""),
			expect:    map[DigestType]string{
				DigestTypeMD5: "d41d8cd98f00b204e9800998ecf8427e",
			},
		},
		{
			name:     "Digests with duplicates",
			digester: NewDigester(DigestTypeMD5, DigestTypeSHA1, DigestTypeSHA256, DigestTypeSHA512, DigestTypeMD5),
			content:  strings.NewReader("this is some content"),
			expect: map[DigestType]string{
				DigestTypeMD5:    "736db904ad222bf88ee6b8d103fceb8e",
				DigestTypeSHA1:   "5ec1a3cb71c75c52cf23934b137985bd2499bd85",
				DigestTypeSHA256: "373993310775a34f5ad48aae265dac65c7abf420dfbaef62819e2cf5aafc64ca",
				DigestTypeSHA512: "47bb28d146567b3be18d06d8468aaa8222183fe6b2a942b17b6a48bbc32bda7213f7dc1acf36677f7710cffa7add3f3656597630bf0d591f34145015f59724e1",
			},
		},
	}

	for _, tc := range testCases {
		got, err := tc.digester.Digest(tc.content)
		if tc.expectErr {
			assert.Error(t, err)
			assert.Equal(t, err, tc.expect)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, got, tc.expect)
		}
	}
}
