package version

var (
	// Git SHA commit set via -ldflags
	BuildCommit string

	// Build date and time in UTC set via -ldflags
	BuildTime string
)
