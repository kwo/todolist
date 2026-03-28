package cli

import (
	"fmt"
	"runtime/debug"
)

// Version is the application semantic version. Release builds should override
// this at link time with -ldflags.
//
//nolint:gochecknoglobals // Set at link time for release builds.
var Version = "dev"

type versionInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit,omitempty"`
	Dirty   bool   `json:"dirty,omitempty"`
	Runtime string `json:"runtime,omitempty"`
}

type versionCommand struct{}

func (c versionCommand) Execute(app *App, options runOptions) error {
	info := versionInfo{Version: Version}

	if app.ReadBuildInfo != nil {
		if buildInfo, ok := app.ReadBuildInfo(); ok {
			applyBuildSettings(&info, buildInfo)
		}
	}

	if options.JSON {
		return writeJSON(app.Stdout, info)
	}

	_, err := fmt.Fprintln(app.Stdout, info.Version)

	return err
}

func applyBuildSettings(info *versionInfo, buildInfo *debug.BuildInfo) {
	info.Runtime = buildInfo.GoVersion

	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			info.Commit = shortCommit(setting.Value)
		case "vcs.modified":
			info.Dirty = setting.Value == "true"
		}
	}
}

func shortCommit(value string) string {
	if len(value) <= 8 {
		return value
	}

	return value[:8]
}
