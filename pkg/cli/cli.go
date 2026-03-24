// Package cli provides the command-line interface for todolist.
package cli

import (
	"fmt"
	"io"
	"os"
	"time"
)

// App is a configurable todolist CLI application.
type App struct {
	// Stdin is the input stream used for reading todo descriptions.
	Stdin io.Reader
	// Stdout is the output stream used for normal command output.
	Stdout io.Writer
	// Stderr is the output stream used for error messages.
	Stderr io.Writer
	// StdinProvided reports whether stdin should be treated as piped input.
	StdinProvided bool
	// TodoDir is the default todo directory used by the application.
	TodoDir string
	// Now returns the current time and is injectable for tests.
	Now func() time.Time
	// LookupEnv resolves environment variables and is injectable for tests.
	LookupEnv func(string) (string, bool)
	// UsageText is the embedded usage documentation printed by the usage command.
	UsageText string
}

type globalOptions struct {
	Directory string
	JSON      bool
	Help      bool
}

type runOptions struct {
	TodoDir string
	JSON    bool
}

type parsedCommand struct {
	name    string
	help    bool
	globals runOptions
	runner  commandRunner
}

type commandRunner interface {
	Execute(*App, runOptions) error
}

// NewApp returns a new CLI application configured with the provided IO streams.
func NewApp(stdin io.Reader, stdout, stderr io.Writer, stdinProvided bool) *App {
	return &App{
		Stdin:         stdin,
		Stdout:        stdout,
		Stderr:        stderr,
		StdinProvided: stdinProvided,
		TodoDir:       "./todo",
		Now:           time.Now,
		LookupEnv:     os.LookupEnv,
	}
}

// Run parses args, executes the selected command, and returns a process exit code.
func (a *App) Run(args []string) int {
	parsed, err := parseArgs(args, a)
	if err != nil {
		_, _ = fmt.Fprintf(a.Stderr, "%s\n", err)

		return 1
	}

	if parsed.help {
		writeHelp(a.Stdout, parsed.name)

		return 0
	}

	if err := parsed.runner.Execute(a, parsed.globals); err != nil {
		_, _ = fmt.Fprintf(a.Stderr, "%s\n", err)

		return 1
	}

	return 0
}

func resolveRunOptions(app *App, globals globalOptions) runOptions {
	todoDir := app.TodoDir
	if app.LookupEnv != nil {
		if value, ok := app.LookupEnv("TODOLIST_DIRECTORY"); ok {
			todoDir = value
		}
	}

	if globals.Directory != "" {
		todoDir = globals.Directory
	}

	return runOptions{
		TodoDir: todoDir,
		JSON:    globals.JSON,
	}
}
