package release

// Taken from: https://golang.org/doc/install/source#environment
// this list needs further curating

// OperatingSystemType enumerates the available operating systems
type OperatingSystemType string

// nolint
const (
	OperatingSystemTypeAndroid   OperatingSystemType = "android"
	OperatingSystemTypeDarwin    OperatingSystemType = "darwin"
	OperatingSystemTypeDragonFly OperatingSystemType = "dragonfly"
	OperatingSystemTypeFreeBSD   OperatingSystemType = "freebsd"
	OperatingSystemTypeLinux     OperatingSystemType = "linux"
	OperatingSystemTypeNetBSD    OperatingSystemType = "netbsd"
	OperatingSystemTypeOpenBSD   OperatingSystemType = "openbsd"
	OperatingSystemTypePlan9     OperatingSystemType = "plan9"
	OperatingSystemTypeSolaris   OperatingSystemType = "solaris"
	OperatingSystemTypeWindows   OperatingSystemType = "windows"
)

// ArchType enumerates the available architectures
type ArchType string

// nolint
const (
	ArchType386      ArchType = "386"   // x86 | x86-32
	ArchTypeamd64    ArchType = "amd64" // x86-64
	ArchTypearm      ArchType = "arm"
	ArchTypearm64    ArchType = "arm64" // AArch64
	ArchTyppc64      ArchType = "ppc64"
	ArchTyppc64le    ArchType = "ppc64le"
	ArchTypemips     ArchType = "mips"
	ArchTypemipsle   ArchType = "mipsle"
	ArchTypemips64   ArchType = "mips64"
	ArchTypemips64le ArchType = "mips64le"
	ArchTypes390x    ArchType = "390x"
)
