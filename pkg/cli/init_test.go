package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kwo/todolist/pkg/cli"
)

func TestInitCreatesTodoDirectoryAndConfig(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	todoDir := filepath.Join(t.TempDir(), "todo")
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = todoDir

	exitCode := app.Run([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("expected init to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}

	info, err := os.Stat(todoDir)
	if err != nil {
		t.Fatalf("stat todo dir: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", todoDir)
	}

	configPath := filepath.Join(todoDir, ".todos")
	//nolint:gosec // Test reads a config file created in a temporary directory.
	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	if string(rawConfig) != "prefix=todo-\n" {
		t.Fatalf("expected default config contents, got %q", string(rawConfig))
	}

	output := stdout.String()
	if !strings.Contains(output, "initialized todo directory: "+todoDir) {
		t.Fatalf("expected initialization output, got %q", output)
	}

	if !strings.Contains(output, "created config file: "+configPath) {
		t.Fatalf("expected config creation output, got %q", output)
	}
}

func TestInitCreatesMissingConfigInExistingDirectory(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	todoDir := filepath.Join(t.TempDir(), "todo")
	if err := os.Mkdir(todoDir, 0o750); err != nil {
		t.Fatalf("create todo dir: %v", err)
	}

	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = todoDir

	exitCode := app.Run([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("expected init to succeed, got %d: %s", exitCode, stderr.String())
	}

	configPath := filepath.Join(todoDir, ".todos")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("stat config: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "todo directory already exists: "+todoDir) {
		t.Fatalf("expected existing directory output, got %q", output)
	}

	if !strings.Contains(output, "created config file: "+configPath) {
		t.Fatalf("expected created config output, got %q", output)
	}
}

func TestInitIsIdempotentAndDoesNotOverwriteConfig(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	todoDir := filepath.Join(t.TempDir(), "todo")
	if err := os.Mkdir(todoDir, 0o750); err != nil {
		t.Fatalf("create todo dir: %v", err)
	}

	configPath := filepath.Join(todoDir, ".todos")
	if err := os.WriteFile(configPath, []byte("prefix=work-\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = todoDir

	exitCode := app.Run([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("expected init to succeed, got %d: %s", exitCode, stderr.String())
	}

	//nolint:gosec // Test reads a config file created in a temporary directory.
	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	if string(rawConfig) != "prefix=work-\n" {
		t.Fatalf("expected existing config to remain unchanged, got %q", string(rawConfig))
	}

	output := stdout.String()
	if !strings.Contains(output, "todo directory already exists: "+todoDir) {
		t.Fatalf("expected existing directory output, got %q", output)
	}

	if !strings.Contains(output, "config file already exists: "+configPath) {
		t.Fatalf("expected existing config output, got %q", output)
	}
}

func TestInitUsesDefaultTodoDirectory(t *testing.T) {
	t.Helper()

	cwd := t.TempDir()
	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}

	if err := os.Chdir(cwd); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}

	t.Cleanup(func() {
		if chdirErr := os.Chdir(previousWD); chdirErr != nil {
			t.Fatalf("restore working directory: %v", chdirErr)
		}
	})

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)

	exitCode := app.Run([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("expected init to succeed, got %d: %s", exitCode, stderr.String())
	}

	if _, err := os.Stat(filepath.Join(cwd, "todo")); err != nil {
		t.Fatalf("stat initialized directory: %v", err)
	}

	if _, err := os.Stat(filepath.Join(cwd, "todo", ".todos")); err != nil {
		t.Fatalf("stat initialized config: %v", err)
	}
}

func TestInitRejectsExtraPositionalArgs(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = filepath.Join(t.TempDir(), "todo")

	exitCode := app.Run([]string{"init", "extra"})
	if exitCode != 1 {
		t.Fatalf("expected init to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "does not accept positional") {
		t.Fatalf("expected positional argument error, got %q", stderr.String())
	}
}

func TestInitRejectsNonDirectoryTodoPath(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	todoPath := filepath.Join(t.TempDir(), "todo")
	if err := os.WriteFile(todoPath, []byte("not a directory\n"), 0o600); err != nil {
		t.Fatalf("write todo path file: %v", err)
	}

	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = todoPath

	exitCode := app.Run([]string{"init"})
	if exitCode != 1 {
		t.Fatalf("expected init to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `is not a directory`) {
		t.Fatalf("expected directory error, got %q", stderr.String())
	}
}

func TestInitRejectsNonRegularConfigPath(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	todoDir := filepath.Join(t.TempDir(), "todo")
	if err := os.Mkdir(todoDir, 0o750); err != nil {
		t.Fatalf("create todo dir: %v", err)
	}

	configPath := filepath.Join(todoDir, ".todos")
	if err := os.Mkdir(configPath, 0o750); err != nil {
		t.Fatalf("create config dir: %v", err)
	}

	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = todoDir

	exitCode := app.Run([]string{"init"})
	if exitCode != 1 {
		t.Fatalf("expected init to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `is not a regular file`) {
		t.Fatalf("expected config path error, got %q", stderr.String())
	}
}
