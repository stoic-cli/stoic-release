package mock

import (
	"github.com/stoic-cli/stoic-release"
	"io/ioutil"
	"strings"
)

var ProjectName = "MyProject"

func ValidDigests() map[release.DigestType]string {
	return map[release.DigestType]string{
		release.DigestTypeMD5:    "736db904ad222bf88ee6b8d103fceb8e",
		release.DigestTypeSHA1:   "5ec1a3cb71c75c52cf23934b137985bd2499bd85",
		release.DigestTypeSHA256: "373993310775a34f5ad48aae265dac65c7abf420dfbaef62819e2cf5aafc64ca",
		release.DigestTypeSHA512: "47bb28d146567b3be18d06d8468aaa8222183fe6b2a942b17b6a48bbc32bda7213f7dc1acf36677f7710cffa7add3f3656597630bf0d591f34145015f59724e1",
	}
}

func ValidArtifacts() []release.Artifact {
	a, _ := release.NewBinaryArtifact(ioutil.NopCloser(strings.NewReader("this is some content")), ProjectName, release.OperatingSystemTypeDarwin, release.ArchTypeamd64)
	a.SetDigests(ValidDigests())
	return []release.Artifact{a}
}
