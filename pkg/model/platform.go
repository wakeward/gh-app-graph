package model

// PlatformAvailability records which GitHub deployment targets can grant a
// permission or enable a toxic combination.
type PlatformAvailability string

const (
	// PlatformAll applies on GitHub.com, Enterprise Cloud, and Enterprise Server.
	PlatformAll PlatformAvailability = "all"
	// PlatformGHESOnly applies only on self-hosted GitHub Enterprise Server.
	PlatformGHESOnly PlatformAvailability = "ghes_only"
)
