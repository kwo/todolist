package cli_test

import (
	"strings"
	"testing"
)

func TestUsagePrintsEmbeddedUsageText(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"usage"})
	if exitCode != 0 {
		t.Fatalf("expected usage to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "## Using the todolist app") {
		t.Fatalf("expected usage heading, got %q", output)
	}

	if !strings.Contains(output, "todolist list --json") {
		t.Fatalf("expected usage content, got %q", output)
	}
}

func TestUsageHelpPrintsCommandUsage(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"usage", "--help"})
	if exitCode != 0 {
		t.Fatalf("expected usage help to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}

	if !strings.Contains(stdout.String(), "todolist usage") {
		t.Fatalf("expected usage command help, got %q", stdout.String())
	}
}
