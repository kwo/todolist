package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kwo/todolist/pkg/cli"
)

func TestDirectoryFlagIsRejected(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)

	exitCode := app.Run([]string{"list", "-d", "./other"})
	if exitCode != 1 {
		t.Fatalf("expected list to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "unknown shorthand flag: 'd' in -d") {
		t.Fatalf("expected unknown flag error, got %q", stderr.String())
	}
}

func TestUnknownCommandFails(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"nope"})
	if exitCode != 1 {
		t.Fatalf("expected unknown command to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `unknown command "nope"`) {
		t.Fatalf("expected unknown command error, got %q", stderr.String())
	}
}
