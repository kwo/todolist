package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kwo/todolist/pkg/cli"
)

func TestDirectoryOptionIsParsedAfterCommand(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = filepath.Join(t.TempDir(), "default-missing")
	app.LookupEnv = func(string) (string, bool) {
		return "", false
	}
	app.Now = func() time.Time {
		return time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC)
	}

	otherDir := filepath.Join(t.TempDir(), "other")
	if err := os.Mkdir(otherDir, 0o750); err != nil {
		t.Fatalf("create other todo dir: %v", err)
	}

	exitCode := app.Run([]string{"add", "-d", otherDir, "--title", "Buy groceries"})
	if exitCode != 0 {
		t.Fatalf("expected add to succeed, got %d: %s", exitCode, stderr.String())
	}

	id := strings.TrimSpace(stdout.String())
	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", "-d", otherDir, id})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	if !strings.Contains(stdout.String(), "title: Buy groceries") {
		t.Fatalf("expected view output from alternate directory, got %q", stdout.String())
	}
}

func TestDirectoryEnvironmentVariableIsUsed(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = filepath.Join(t.TempDir(), "default-missing")
	envDir := filepath.Join(t.TempDir(), "env-todo")
	app.LookupEnv = func(key string) (string, bool) {
		if key != "TODOLIST_DIRECTORY" {
			return "", false
		}

		return envDir, true
	}
	app.Now = func() time.Time {
		return time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC)
	}

	if err := os.Mkdir(envDir, 0o750); err != nil {
		t.Fatalf("create env todo dir: %v", err)
	}

	exitCode := app.Run([]string{"add", "--title", "Buy groceries"})
	if exitCode != 0 {
		t.Fatalf("expected add to succeed, got %d: %s", exitCode, stderr.String())
	}

	id := strings.TrimSpace(stdout.String())
	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", id})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	if !strings.Contains(stdout.String(), "title: Buy groceries") {
		t.Fatalf("expected view output from env directory, got %q", stdout.String())
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
