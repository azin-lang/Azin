package platform

type OS uint8

const (
	OSUnknown OS = iota

	// Mainstream
	Linux
	Windows
	Darwin
	Android

	// BSDs
	FreeBSD
	OpenBSD
	NetBSD
	DragonFlyBSD

	// Solaris family
	Solaris
	Illumos

	// Other Unix
	AIX

	// WebAssembly
	WASI
)
