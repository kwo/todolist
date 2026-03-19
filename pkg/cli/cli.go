// Package cli provides the command-line interface for tasklist.
package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"

	"github.com/kwo/tasklist/pkg/tasklist"
)

// App is a configurable tasklist CLI application.
type App struct {
	// Stdin is the input stream used for reading task descriptions.
	Stdin io.Reader
	// Stdout is the output stream used for normal command output.
	Stdout io.Writer
	// Stderr is the output stream used for error messages.
	Stderr io.Writer
	// StdinProvided reports whether stdin should be treated as piped input.
	StdinProvided bool
	// TaskDir is the task directory used by the application.
	TaskDir string
	// Now returns the current time and is injectable for tests.
	Now func() time.Time
}

type rootCommand struct{}

type addCommand struct {
	app  *App
	Args struct {
		Title string `positional-arg-name:"title" required:"yes"`
	} `positional-args:"yes"`
}

type deleteCommand struct {
	app  *App
	Args struct {
		Task string `positional-arg-name:"task" required:"yes"`
	} `positional-args:"yes"`
}

type listCommand struct {
	app *App
}

type updateCommand struct {
	app   *App
	Title string `long:"title" description:"replace the task title"`
	Args  struct {
		Task string `positional-arg-name:"task" required:"yes"`
	} `positional-args:"yes"`
}

type viewCommand struct {
	app  *App
	Args struct {
		Task string `positional-arg-name:"task" required:"yes"`
	} `positional-args:"yes"`
}

// NewApp returns a new CLI application configured with the provided IO streams.
func NewApp(stdin io.Reader, stdout, stderr io.Writer, stdinProvided bool) *App {
	return &App{
		Stdin:         stdin,
		Stdout:        stdout,
		Stderr:        stderr,
		StdinProvided: stdinProvided,
		TaskDir:       "./tasks",
		Now:           time.Now,
	}
}

// Run parses args, executes the selected command, and returns a process exit code.
func (a *App) Run(args []string) int {
	parser := newParser(a)
	_, err := parser.ParseArgs(args)
	if err == nil {
		return 0
	}

	if isHelp(err) {
		parser.WriteHelp(a.Stdout)
		return 0
	}

	_, _ = fmt.Fprintf(a.Stderr, "%s\n", err)

	return 1
}

func newParser(app *App) *flags.Parser {
	parser := flags.NewParser(&rootCommand{}, flags.HelpFlag)
	parser.Name = "tasklist"

	_, _ = parser.AddCommand("add", "Create a task", "Create a task", &addCommand{app: app})
	_, _ = parser.AddCommand("list", "List tasks", "List tasks", &listCommand{app: app})
	_, _ = parser.AddCommand("view", "View a task", "View a task", &viewCommand{app: app})
	_, _ = parser.AddCommand("update", "Update a task", "Update a task", &updateCommand{app: app})
	_, _ = parser.AddCommand("delete", "Delete a task", "Delete a task", &deleteCommand{app: app})

	return parser
}

func (c *addCommand) Execute(args []string) error {
	if err := rejectExtras(args); err != nil {
		return err
	}

	title := strings.TrimSpace(c.Args.Title)
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	description, err := readDescription(c.app.Stdin, c.app.StdinProvided)
	if err != nil {
		return err
	}

	now := tasklist.NormalizeTimestamp(c.app.Now())
	store := tasklist.NewStore(c.app.TaskDir)
	value := tasklist.Task{
		Title:        title,
		CreatedAt:    now,
		LastModified: now,
		Description:  description,
	}
	value.ID = tasklist.GenerateID(value, store.Exists)

	if err := store.Create(value); err != nil {
		return err
	}

	_, err = fmt.Fprintf(c.app.Stdout, "%s\n", value.ID)

	return err
}

func (c *listCommand) Execute(args []string) error {
	if err := rejectExtras(args); err != nil {
		return err
	}

	tasks, err := tasklist.NewStore(c.app.TaskDir).List()
	if err != nil {
		return err
	}

	for _, value := range tasks {
		if _, err = fmt.Fprintf(c.app.Stdout, "%s\t%s\n", value.ID, value.Title); err != nil {
			return err
		}
	}

	return nil
}

func (c *viewCommand) Execute(args []string) error {
	if err := rejectExtras(args); err != nil {
		return err
	}

	raw, err := tasklist.NewStore(c.app.TaskDir).GetRaw(c.Args.Task)
	if err != nil {
		return err
	}

	_, err = c.app.Stdout.Write(raw)

	return err
}

func (c *updateCommand) Execute(args []string) error {
	if err := rejectExtras(args); err != nil {
		return err
	}

	titleProvided := c.Title != ""
	description, descriptionProvided, err := readOptionalDescription(c.app.Stdin, c.app.StdinProvided)
	if err != nil {
		return err
	}

	if !titleProvided && !descriptionProvided {
		return fmt.Errorf("update requires --title or stdin description input")
	}

	store := tasklist.NewStore(c.app.TaskDir)
	value, err := store.Get(c.Args.Task)
	if err != nil {
		return err
	}

	if titleProvided {
		value.Title = strings.TrimSpace(c.Title)
		if value.Title == "" {
			return fmt.Errorf("title cannot be empty")
		}
	}

	if descriptionProvided {
		value.Description = description
	}

	value.LastModified = tasklist.NormalizeTimestamp(c.app.Now())

	return store.Update(value)
}

func (c *deleteCommand) Execute(args []string) error {
	if err := rejectExtras(args); err != nil {
		return err
	}

	return tasklist.NewStore(c.app.TaskDir).Delete(c.Args.Task)
}

func readDescription(reader io.Reader, provided bool) (string, error) {
	value, _, err := readOptionalDescription(reader, provided)

	return value, err
}

func readOptionalDescription(reader io.Reader, provided bool) (string, bool, error) {
	if !provided {
		return "", false, nil
	}

	raw, err := io.ReadAll(reader)
	if err != nil {
		return "", false, fmt.Errorf("read stdin: %w", err)
	}

	return string(raw), true, nil
}

func rejectExtras(args []string) error {
	if len(args) == 0 {
		return nil
	}

	return fmt.Errorf("unknown arguments: %s", strings.Join(args, " "))
}

func isHelp(err error) bool {
	typed := &flags.Error{}
	if !asFlagsError(err, typed) {
		return false
	}

	return typed.Type == flags.ErrHelp
}

func asFlagsError(err error, target *flags.Error) bool {
	typed := &flags.Error{}
	ok := errors.As(err, &typed)
	if !ok {
		return false
	}

	*target = *typed

	return true
}
