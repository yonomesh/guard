package cmd

import (
	"fmt"
	"strings"
)

func init() {
	// do nothing
}

// Flags that control program flow or startup
var (
	conf    string
	version bool
	plugins bool

	// LogFlags are initially set to 0 for no extra output
	LogFlags int
)

// usage
// go build -ldflags "-X 'main.appVersion=1.2.3' -X 'main.devBuild=false' -X 'main.buildDate=$(date -u)'"
var (
	appVersion       = "(untracked dev build)" // inferred at startup
	devBuild         = true                    // inferred at startup
	buildDate        string                    // date -u
	gitTag           string                    // git describe --exact-match HEAD 2> /dev/null
	gitNearestTag    string                    // git describe --abbrev=0 --tags HEAD
	gitShortStat     string                    // git diff-index --shortstat
	gitFilesModified string                    // git diff-index --name-only HEAD

	// Gitcommit contains the commit where we built Guard from.
	GitCommit string
)

// setVersion get the version information
// based on variables set by -ldflags.
func setVersion() {
	// A development build is one that's not at a tag or has uncommitted changes
	devBuild = gitTag == "" || gitShortStat != ""

	// Only set the appVersion if -ldflags was used
	if gitNearestTag != "" || gitTag != "" {
		if devBuild && gitNearestTag != "" {
			appVersion = fmt.Sprintf("%s (+%s %s)", strings.TrimPrefix(gitNearestTag, "v"), GitCommit, buildDate)
		} else if gitTag != "" {
			appVersion = strings.TrimPrefix(gitTag, "v")
		}
	}
}

func Main() {
	fmt.Println("Test")
}
