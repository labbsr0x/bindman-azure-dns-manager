package version

import (
	"fmt"
	"runtime"
)

// Default build-time variable.
// These values are overridden via ldflags
var (
	Version   = "unknown-version"
	GitCommit = "unknown-commit"
	BuildTime = "unknown-buildtime"

	Revision  string
	Branch    string
	GoVersion = runtime.Version()
)

const versionF = `Bindman Azure DNS Manager
  Version: %s
  GitCommit: %s
  BuildTime: %s
`

func FormattedMessage() string {
	return fmt.Sprintf(versionF, Version, GitCommit, BuildTime)
}
