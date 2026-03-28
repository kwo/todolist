package cli_test

import (
	"encoding/json"
	"runtime/debug"
	"testing"

	"github.com/kwo/todolist/pkg/cli"
)

func TestVersionPrintsDefaultVersion(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	app.ReadBuildInfo = func() (*debug.BuildInfo, bool) {
		return nil, false
	}

	originalVersion := cli.Version
	cli.Version = "dev"
	t.Cleanup(func() {
		cli.Version = originalVersion
	})

	exitCode := app.Run([]string{"version"})
	if exitCode != 0 {
		t.Fatalf("expected version to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stdout.String() != "dev\n" {
		t.Fatalf("unexpected version output %q", stdout.String())
	}
}

func TestVersionTextOutputPrintsOnlyVersion(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	app.ReadBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{Settings: []debug.BuildSetting{
			{Key: "vcs.revision", Value: "1234567890abcdef"},
			{Key: "vcs.modified", Value: "true"},
		}}, true
	}

	originalVersion := cli.Version
	cli.Version = "v1.2.3"
	t.Cleanup(func() {
		cli.Version = originalVersion
	})

	exitCode := app.Run([]string{"version"})
	if exitCode != 0 {
		t.Fatalf("expected version to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stdout.String() != "v1.2.3\n" {
		t.Fatalf("unexpected version output %q", stdout.String())
	}
}

func TestVersionJSONIncludesAvailableFields(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	app.ReadBuildInfo = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			GoVersion: "go1.26.1",
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "deadbeefcafebabe"},
				{Key: "vcs.modified", Value: "false"},
			},
		}, true
	}

	originalVersion := cli.Version
	cli.Version = "v9.9.9"
	t.Cleanup(func() {
		cli.Version = originalVersion
	})

	exitCode := app.Run([]string{"version", "--json"})
	if exitCode != 0 {
		t.Fatalf("expected json version to succeed, got %d: %s", exitCode, stderr.String())
	}

	var result map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal version json: %v; output=%q", err, stdout.String())
	}

	if result["version"] != "v9.9.9" {
		t.Fatalf("expected version in json, got %+v", result)
	}

	if result["runtime"] != "go1.26.1" {
		t.Fatalf("expected runtime in json, got %+v", result)
	}

	if result["commit"] != "deadbeef" {
		t.Fatalf("expected shortened commit in json, got %+v", result)
	}

	if _, ok := result["dirty"]; ok {
		t.Fatalf("expected dirty to be omitted when false, got %+v", result)
	}
}
