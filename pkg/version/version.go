package version

import (
	"encoding/json"
	"fmt"
	"runtime"
)

var (
	// GitVersion is semantic version.
	GitVersion = "v0.0.0-master+$Format:%h$"
	// BuildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ').
	BuildDate = "1970-01-01T00:00:00Z"
	// GitCommit sha1 from git, output of $(git rev-parse HEAD).
	GitCommit = "$Format:%H$"
	// GitTreeState state of git tree, either "clean" or "dirty".
	GitTreeState = ""
)

// Info contains versioning information.
type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// Get returns the overall codebase version.
func Get() Info {
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a human-readable version string.
func (i Info) String() string {
	return i.GitVersion
}

// ToJson returns the JSON format of version info.
func (i Info) ToJson() string {
	s, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(s)
}
