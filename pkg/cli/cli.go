// Package cli provides the command-line interface for todolist.
package cli

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"

	"github.com/kwo/todolist/pkg/todolist"
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
	// TodoDir is the todo directory used by the application.
	TodoDir string
	// Now returns the current time and is injectable for tests.
	Now func() time.Time
}

type rootCommand struct{}

type addCommand struct {
	app      *App
	Status   string `short:"s" long:"status" description:"set the todo status"`
	Priority int    `short:"p" long:"priority" description:"set the todo priority"`
	Args     struct {
		Title string `positional-arg-name:"title" required:"yes"`
	} `positional-args:"yes"`
}

type deleteCommand struct {
	app  *App
	Args struct {
		Todo string `positional-arg-name:"todo" required:"yes"`
	} `positional-args:"yes"`
}

type listCommand struct {
	app      *App
	Status   string `short:"s" long:"status" description:"filter by status; prefix with ! to exclude"`
	Priority string `short:"p" long:"priority" description:"filter by priority; supports 3, !3, >3, <3"`
}

type updateCommand struct {
	app      *App
	Title    string `long:"title" description:"replace the todo title"`
	Status   string `short:"s" long:"status" description:"replace the todo status"`
	Priority int    `short:"p" long:"priority" description:"replace the todo priority"`
	Args     struct {
		Todo string `positional-arg-name:"todo" required:"yes"`
	} `positional-args:"yes"`
}

type viewCommand struct {
	app  *App
	Args struct {
		Todo string `positional-arg-name:"todo" required:"yes"`
	} `positional-args:"yes"`
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
	parser.Name = "todolist"

	_, _ = parser.AddCommand("add", "Create a todo", "Create a todo", &addCommand{app: app})
	_, _ = parser.AddCommand("list", "List todos", "List todos", &listCommand{app: app})
	_, _ = parser.AddCommand("view", "View a todo", "View a todo", &viewCommand{app: app})
	_, _ = parser.AddCommand("update", "Update a todo", "Update a todo", &updateCommand{app: app})
	_, _ = parser.AddCommand("delete", "Delete a todo", "Delete a todo", &deleteCommand{app: app})

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

	status := strings.TrimSpace(c.Status)
	if status == "" {
		status = todolist.DefaultStatus
	}

	if err := todolist.ValidateStatus(status); err != nil {
		return err
	}

	priority := c.Priority
	if priority == 0 {
		priority = todolist.DefaultPriority
	}

	if err := todolist.ValidatePriority(priority); err != nil {
		return err
	}

	now := todolist.NormalizeTimestamp(c.app.Now())
	store := todolist.NewStore(c.app.TodoDir)
	value := todolist.Todo{
		Title:        title,
		Status:       status,
		Priority:     priority,
		CreatedAt:    now,
		LastModified: now,
		Description:  description,
	}
	value.ID = todolist.GenerateID(value, store.Exists)

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

	statusFilter, excludeStatus, err := parseStatusFilter(c.Status)
	if err != nil {
		return err
	}

	priorityFilter, err := parsePriorityFilter(c.Priority)
	if err != nil {
		return err
	}

	todos, err := todolist.NewStore(c.app.TodoDir).List()
	if err != nil {
		return err
	}

	for _, value := range todos {
		if statusFilter != "" {
			matchesStatus := value.Status == statusFilter
			if (!excludeStatus && !matchesStatus) || (excludeStatus && matchesStatus) {
				continue
			}
		}

		if priorityFilter != nil && !priorityFilter(value.Priority) {
			continue
		}

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

	raw, err := todolist.NewStore(c.app.TodoDir).GetRaw(c.Args.Todo)
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

	description, descriptionProvided, err := readOptionalDescription(c.app.Stdin, c.app.StdinProvided)
	if err != nil {
		return err
	}

	store := todolist.NewStore(c.app.TodoDir)
	value, err := store.Get(c.Args.Todo)
	if err != nil {
		return err
	}

	titleProvided := c.Title != ""
	statusProvided := c.Status != ""
	priorityProvided := c.Priority != 0

	if titleProvided {
		value.Title = strings.TrimSpace(c.Title)
		if value.Title == "" {
			return fmt.Errorf("title cannot be empty")
		}
	}

	if statusProvided {
		value.Status = strings.TrimSpace(c.Status)
		if err := todolist.ValidateStatus(value.Status); err != nil {
			return err
		}
	}

	if priorityProvided {
		if err := todolist.ValidatePriority(c.Priority); err != nil {
			return err
		}

		value.Priority = c.Priority
	}

	if descriptionProvided {
		value.Description = description
	}

	if !titleProvided && !statusProvided && !priorityProvided && !descriptionProvided {
		return fmt.Errorf("update requires --title, --status, --priority, or stdin description input")
	}

	value.LastModified = todolist.NormalizeTimestamp(c.app.Now())

	return store.Update(value)
}

func (c *deleteCommand) Execute(args []string) error {
	if err := rejectExtras(args); err != nil {
		return err
	}

	return todolist.NewStore(c.app.TodoDir).Delete(c.Args.Todo)
}

func readDescription(reader io.Reader, provided bool) (string, error) {
	value, _, err := readOptionalDescription(reader, provided)

	return value, err
}

func parseStatusFilter(raw string) (string, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, nil
	}

	exclude := strings.HasPrefix(value, "!")
	if exclude {
		value = strings.TrimSpace(strings.TrimPrefix(value, "!"))
	}

	if err := todolist.ValidateStatus(value); err != nil {
		return "", false, err
	}

	return value, exclude, nil
}

func parsePriorityFilter(raw string) (func(int) bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, nil
	}

	operator := "="
	if strings.HasPrefix(value, ".") || strings.HasPrefix(value, "+") || strings.HasPrefix(value, "-") {
		operator = value[:1]
		value = strings.TrimSpace(value[1:])
	}

	priority, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid priority filter %q: must use n, .n, +n, or -n", raw)
	}

	if priority < 0 || priority > todolist.DefaultPriority {
		return nil, fmt.Errorf("invalid priority filter %q: priority must be between 0 and %d", raw, todolist.DefaultPriority)
	}

	switch operator {
	case "=":
		return func(candidate int) bool { return candidate == priority }, nil
	case ".":
		return func(candidate int) bool { return candidate != priority }, nil
	case "+":
		return func(candidate int) bool { return candidate > priority }, nil
	case "-":
		return func(candidate int) bool { return candidate < priority }, nil
	default:
		return nil, fmt.Errorf("invalid priority filter %q", raw)
	}
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
