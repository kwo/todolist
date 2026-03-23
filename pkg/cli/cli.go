// Package cli provides the command-line interface for todolist.
package cli

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

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
	// TodoDir is the default todo directory used by the application.
	TodoDir string
	// Now returns the current time and is injectable for tests.
	Now func() time.Time
	// LookupEnv resolves environment variables and is injectable for tests.
	LookupEnv func(string) (string, bool)
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

type addCommand struct {
	Title    string
	Status   string
	Priority int
}

type deleteCommand struct {
	Todo string
}

type listCommand struct {
	StatusFilter   string
	ExcludeStatus  bool
	PriorityFilter string
}

type updateCommand struct {
	Todo             string
	Title            string
	Status           string
	Priority         int
	TitleProvided    bool
	StatusProvided   bool
	PriorityProvided bool
}

type viewCommand struct {
	Todo string
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

func parseArgs(args []string, app *App) (parsedCommand, error) {
	if len(args) == 0 {
		return parsedCommand{}, fmt.Errorf("missing command")
	}

	if isHelpToken(args[0]) {
		return parsedCommand{name: "", help: true, globals: resolveRunOptions(app, globalOptions{})}, nil
	}

	commandName := args[0]
	commandValues, globals, err := parseGlobalOptions(args[1:])
	if err != nil {
		return parsedCommand{}, err
	}

	parsed := parsedCommand{
		name:    commandName,
		help:    globals.Help,
		globals: resolveRunOptions(app, globals),
	}

	switch commandName {
	case "add":
		if parsed.help {
			return parsed, nil
		}

		runner, addErr := parseAddCommand(commandValues)
		if addErr != nil {
			return parsedCommand{}, addErr
		}

		parsed.runner = runner
	case "list":
		if parsed.help {
			return parsed, nil
		}

		runner, listErr := parseListCommand(commandValues)
		if listErr != nil {
			return parsedCommand{}, listErr
		}

		parsed.runner = runner
	case "view":
		if parsed.help {
			return parsed, nil
		}

		runner, viewErr := parseSingleTodoCommand("view", commandValues)
		if viewErr != nil {
			return parsedCommand{}, viewErr
		}

		parsed.runner = runner
	case "update":
		if parsed.help {
			return parsed, nil
		}

		runner, updateErr := parseUpdateCommand(commandValues)
		if updateErr != nil {
			return parsedCommand{}, updateErr
		}

		parsed.runner = runner
	case "delete":
		if parsed.help {
			return parsed, nil
		}

		runner, deleteErr := parseSingleTodoCommand("delete", commandValues)
		if deleteErr != nil {
			return parsedCommand{}, deleteErr
		}

		parsed.runner = deleteCommand(runner)
	default:
		return parsedCommand{}, fmt.Errorf("unknown command %q", commandName)
	}

	return parsed, nil
}

func parseGlobalOptions(args []string) ([]string, globalOptions, error) {
	values := make([]string, 0, len(args))
	options := globalOptions{}

	for index := 0; index < len(args); index++ {
		arg := args[index]

		switch {
		case arg == "-h" || arg == "--help":
			options.Help = true
		case arg == "--json":
			options.JSON = true
		case arg == "-d" || arg == "--directory":
			if index+1 >= len(args) {
				return nil, globalOptions{}, fmt.Errorf("%s requires a directory value", arg)
			}

			index++
			options.Directory = args[index]
		case strings.HasPrefix(arg, "-d="):
			options.Directory = strings.TrimPrefix(arg, "-d=")
		case strings.HasPrefix(arg, "--directory="):
			options.Directory = strings.TrimPrefix(arg, "--directory=")
		default:
			values = append(values, arg)
		}
	}

	return values, options, nil
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

func parseAddCommand(values []string) (addCommand, error) {
	command := addCommand{
		Status:   todolist.DefaultStatus,
		Priority: todolist.DefaultPriority,
	}

	titleAssigned := false
	statusAssigned := false
	priorityAssigned := false

	for _, raw := range values {
		key, value, explicit := parseAssignment(raw)
		if explicit {
			switch key {
			case "title":
				if titleAssigned {
					return addCommand{}, fmt.Errorf("title was provided more than once")
				}

				command.Title = value
				titleAssigned = true
			case "status":
				if statusAssigned {
					return addCommand{}, fmt.Errorf("status was provided more than once")
				}

				if err := todolist.ValidateStatus(strings.TrimSpace(value)); err != nil {
					return addCommand{}, err
				}

				command.Status = strings.TrimSpace(value)
				statusAssigned = true
			case "priority":
				if priorityAssigned {
					return addCommand{}, fmt.Errorf("priority was provided more than once")
				}

				priority, err := parseExplicitPriority(value)
				if err != nil {
					return addCommand{}, err
				}

				command.Priority = priority
				priorityAssigned = true
			default:
				return addCommand{}, fmt.Errorf("unknown add field %q", key)
			}

			continue
		}

		if !statusAssigned {
			status, ok := recognizeStatusValue(raw)
			if ok {
				command.Status = status
				statusAssigned = true

				continue
			}
		}

		if !priorityAssigned {
			priority, ok := recognizePriorityValue(raw)
			if ok {
				command.Priority = priority
				priorityAssigned = true

				continue
			}
		}

		if !titleAssigned {
			command.Title = raw
			titleAssigned = true

			continue
		}

		return addCommand{}, fmt.Errorf("cannot assign value %q", raw)
	}

	if !titleAssigned {
		return addCommand{}, fmt.Errorf("add requires a title")
	}

	return command, nil
}

func parseListCommand(values []string) (listCommand, error) {
	command := listCommand{}
	statusAssigned := false
	priorityAssigned := false

	for _, raw := range values {
		key, value, explicit := parseAssignment(raw)
		if explicit {
			switch key {
			case "status":
				if statusAssigned {
					return listCommand{}, fmt.Errorf("status filter was provided more than once")
				}

				statusFilter, excludeStatus, err := parseExplicitStatusFilter(value)
				if err != nil {
					return listCommand{}, err
				}

				command.StatusFilter = statusFilter
				command.ExcludeStatus = excludeStatus
				statusAssigned = true
			case "priority":
				if priorityAssigned {
					return listCommand{}, fmt.Errorf("priority filter was provided more than once")
				}

				if _, err := parsePriorityFilter(value); err != nil {
					return listCommand{}, err
				}

				command.PriorityFilter = strings.TrimSpace(value)
				priorityAssigned = true
			default:
				return listCommand{}, fmt.Errorf("unknown list field %q", key)
			}

			continue
		}

		if !statusAssigned {
			statusFilter, excludeStatus, ok, err := recognizeStatusFilter(raw)
			if err != nil {
				return listCommand{}, err
			}

			if ok {
				command.StatusFilter = statusFilter
				command.ExcludeStatus = excludeStatus
				statusAssigned = true

				continue
			}
		}

		if !priorityAssigned {
			priorityFilter, ok, err := recognizePriorityFilter(raw)
			if err != nil {
				return listCommand{}, err
			}

			if ok {
				command.PriorityFilter = priorityFilter
				priorityAssigned = true

				continue
			}
		}

		return listCommand{}, fmt.Errorf("cannot assign value %q", raw)
	}

	return command, nil
}

func parseSingleTodoCommand(name string, values []string) (viewCommand, error) {
	if len(values) == 0 {
		return viewCommand{}, fmt.Errorf("%s requires a todo id", name)
	}

	if len(values) > 1 {
		return viewCommand{}, fmt.Errorf("cannot assign value %q", values[1])
	}

	return viewCommand{Todo: values[0]}, nil
}

func parseUpdateCommand(values []string) (updateCommand, error) {
	if len(values) == 0 {
		return updateCommand{}, fmt.Errorf("update requires a todo id")
	}

	command := updateCommand{Todo: values[0]}
	titleAssigned := false
	statusAssigned := false
	priorityAssigned := false

	for _, raw := range values[1:] {
		key, value, explicit := parseAssignment(raw)
		if explicit {
			switch key {
			case "title":
				if titleAssigned {
					return updateCommand{}, fmt.Errorf("title was provided more than once")
				}

				command.Title = value
				command.TitleProvided = true
				titleAssigned = true
			case "status":
				if statusAssigned {
					return updateCommand{}, fmt.Errorf("status was provided more than once")
				}

				if err := todolist.ValidateStatus(strings.TrimSpace(value)); err != nil {
					return updateCommand{}, err
				}

				command.Status = strings.TrimSpace(value)
				command.StatusProvided = true
				statusAssigned = true
			case "priority":
				if priorityAssigned {
					return updateCommand{}, fmt.Errorf("priority was provided more than once")
				}

				priority, err := parseExplicitPriority(value)
				if err != nil {
					return updateCommand{}, err
				}

				command.Priority = priority
				command.PriorityProvided = true
				priorityAssigned = true
			default:
				return updateCommand{}, fmt.Errorf("unknown update field %q", key)
			}

			continue
		}

		if !statusAssigned {
			status, ok := recognizeStatusValue(raw)
			if ok {
				command.Status = status
				command.StatusProvided = true
				statusAssigned = true

				continue
			}
		}

		if !priorityAssigned {
			priority, ok := recognizePriorityValue(raw)
			if ok {
				command.Priority = priority
				command.PriorityProvided = true
				priorityAssigned = true

				continue
			}
		}

		if !titleAssigned {
			command.Title = raw
			command.TitleProvided = true
			titleAssigned = true

			continue
		}

		return updateCommand{}, fmt.Errorf("cannot assign value %q", raw)
	}

	return command, nil
}

func (c addCommand) Execute(app *App, options runOptions) error {
	title := strings.TrimSpace(c.Title)
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	description, err := readDescription(app.Stdin, app.StdinProvided)
	if err != nil {
		return err
	}

	now := todolist.NormalizeTimestamp(app.Now())
	store := todolist.NewStore(options.TodoDir)
	config, err := todolist.LoadConfig(options.TodoDir)
	if err != nil {
		return err
	}

	value := todolist.Todo{
		Title:        title,
		Status:       c.Status,
		Priority:     c.Priority,
		CreatedAt:    now,
		LastModified: now,
		Description:  description,
	}
	value.ID = todolist.GenerateIDWithPrefix(value, config.Prefix, store.Exists)

	if err := store.Create(value); err != nil {
		return err
	}

	_, err = fmt.Fprintf(app.Stdout, "%s\n", value.ID)

	return err
}

func (c listCommand) Execute(app *App, options runOptions) error {
	priorityFilter, err := parsePriorityFilter(c.PriorityFilter)
	if err != nil {
		return err
	}

	todos, err := todolist.NewStore(options.TodoDir).List()
	if err != nil {
		return err
	}

	for _, value := range todos {
		if c.StatusFilter != "" {
			matchesStatus := value.Status == c.StatusFilter
			if (!c.ExcludeStatus && !matchesStatus) || (c.ExcludeStatus && matchesStatus) {
				continue
			}
		}

		if priorityFilter != nil && !priorityFilter(value.Priority) {
			continue
		}

		if _, err = fmt.Fprintf(app.Stdout, "%s\t%s\n", value.ID, value.Title); err != nil {
			return err
		}
	}

	return nil
}

func (c viewCommand) Execute(app *App, options runOptions) error {
	raw, err := todolist.NewStore(options.TodoDir).GetRaw(c.Todo)
	if err != nil {
		return err
	}

	_, err = app.Stdout.Write(raw)

	return err
}

func (c updateCommand) Execute(app *App, options runOptions) error {
	description, descriptionProvided, err := readOptionalDescription(app.Stdin, app.StdinProvided)
	if err != nil {
		return err
	}

	store := todolist.NewStore(options.TodoDir)
	value, err := store.Get(c.Todo)
	if err != nil {
		return err
	}

	if c.TitleProvided {
		value.Title = strings.TrimSpace(c.Title)
		if value.Title == "" {
			return fmt.Errorf("title cannot be empty")
		}
	}

	if c.StatusProvided {
		value.Status = c.Status
	}

	if c.PriorityProvided {
		value.Priority = c.Priority
	}

	if descriptionProvided {
		value.Description = description
	}

	if !c.TitleProvided && !c.StatusProvided && !c.PriorityProvided && !descriptionProvided {
		return fmt.Errorf("update requires a title, status, priority, or stdin description input")
	}

	value.LastModified = todolist.NormalizeTimestamp(app.Now())

	return store.Update(value)
}

func (c deleteCommand) Execute(app *App, options runOptions) error {
	return todolist.NewStore(options.TodoDir).Delete(c.Todo)
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

	exclude := strings.HasSuffix(value, "!")
	if exclude {
		value = strings.TrimSpace(strings.TrimSuffix(value, "!"))
	}

	if err := todolist.ValidateStatus(value); err != nil {
		return "", false, err
	}

	return value, exclude, nil
}

func parseExplicitStatusFilter(raw string) (string, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, fmt.Errorf("invalid status %q: must be one of todo, wip, done", value)
	}

	return parseStatusFilter(value)
}

func parsePriorityFilter(raw string) (func(int) bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, nil
	}

	operator := "="
	switch {
	case hasPriorityFilterPrefix(value):
		operator = value[:1]
		value = strings.TrimSpace(value[1:])
	case hasPriorityFilterSuffix(value):
		operator = value[len(value)-1:]
		value = strings.TrimSpace(value[:len(value)-1])
	}

	priority, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid priority filter %q: must use n, n!, n+, or n-", raw)
	}

	if priority < 0 || priority > todolist.DefaultPriority {
		return nil, fmt.Errorf("invalid priority filter %q: priority must be between 0 and %d", raw, todolist.DefaultPriority)
	}

	switch operator {
	case "=":
		return func(candidate int) bool { return candidate == priority }, nil
	case ".", "!":
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

func parseAssignment(raw string) (string, string, bool) {
	index := strings.Index(raw, "=")
	if index <= 0 {
		return "", "", false
	}

	return strings.TrimSpace(raw[:index]), raw[index+1:], true
}

func parseExplicitPriority(raw string) (int, error) {
	value := strings.TrimSpace(raw)
	priority, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid priority %q: must be between 1 and 5", raw)
	}

	if err := todolist.ValidatePriority(priority); err != nil {
		return 0, err
	}

	return priority, nil
}

func recognizeStatusValue(raw string) (string, bool) {
	value := strings.TrimSpace(raw)
	if err := todolist.ValidateStatus(value); err != nil {
		return "", false
	}

	return value, true
}

func recognizePriorityValue(raw string) (int, bool) {
	value := strings.TrimSpace(raw)
	if value == "" || !isDigits(value) {
		return 0, false
	}

	priority, err := strconv.Atoi(value)
	if err != nil {
		return 0, false
	}

	if err := todolist.ValidatePriority(priority); err != nil {
		return 0, false
	}

	return priority, true
}

func recognizeStatusFilter(raw string) (string, bool, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, false, nil
	}

	if strings.HasSuffix(value, "!") {
		status, exclude, err := parseStatusFilter(value)
		if err == nil {
			return status, exclude, true, nil
		}
	}

	if strings.HasPrefix(value, "!") {
		status, exclude, err := parseStatusFilter(value)

		return status, exclude, true, err
	}

	status, ok := recognizeStatusValue(value)
	if !ok {
		return "", false, false, nil
	}

	return status, false, true, nil
}

func recognizePriorityFilter(raw string) (string, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, nil
	}

	if !isPriorityFilterShape(value) {
		return "", false, nil
	}

	if _, err := parsePriorityFilter(value); err != nil {
		return "", true, err
	}

	return value, true, nil
}

func isPriorityFilterShape(value string) bool {
	if isDigits(value) {
		return true
	}

	if len(value) < 2 {
		return value == "." || value == "+" || value == "-" || value == "!"
	}

	if hasPriorityFilterPrefix(value) {
		return isDigits(strings.TrimSpace(value[1:])) || strings.TrimSpace(value[1:]) == ""
	}

	if hasPriorityFilterSuffix(value) {
		return isDigits(strings.TrimSpace(value[:len(value)-1])) || strings.TrimSpace(value[:len(value)-1]) == ""
	}

	return false
}

func hasPriorityFilterPrefix(value string) bool {
	switch value[0] {
	case '.', '+', '-':
		return true
	default:
		return false
	}
}

func hasPriorityFilterSuffix(value string) bool {
	switch value[len(value)-1] {
	case '!', '+', '-':
		return true
	default:
		return false
	}
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}

	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func isHelpToken(value string) bool {
	return value == "-h" || value == "--help"
}

func writeHelp(writer io.Writer, command string) {
	if command == "" {
		_, _ = fmt.Fprint(writer, rootHelp())

		return
	}

	_, _ = fmt.Fprint(writer, commandHelp(command))
}

func rootHelp() string {
	return strings.Join([]string{
		"Usage:",
		"  todolist <command> [global options] [command values...]",
		"",
		"Commands:",
		"  add <title> [<status>] [<priority>]",
		"  list [<status-filter>] [<priority-filter>]",
		"  view <todo>",
		"  update <todo> [<title>] [<status>] [<priority>]",
		"  delete <todo>",
		"",
		"Global options:",
		"  -d, --directory <dir>  use a specific todo directory",
		"      --json             enable JSON output",
		"  -h, --help             show help",
		"",
	}, "\n")
}

func commandHelp(command string) string {
	switch command {
	case "add":
		return strings.Join([]string{
			"Usage:",
			"  todolist add [global options] <title> [<status>] [<priority>]",
			"  todolist add [global options] title=<title> [status=<status>] [priority=<priority>]",
			"",
		}, "\n")
	case "list":
		return strings.Join([]string{
			"Usage:",
			"  todolist list [global options] [<status-filter>] [<priority-filter>]",
			"  todolist list [global options] [status=<status-filter>] [priority=<priority-filter>]",
			"",
		}, "\n")
	case "view":
		return strings.Join([]string{
			"Usage:",
			"  todolist view [global options] <todo>",
			"",
		}, "\n")
	case "update":
		return strings.Join([]string{
			"Usage:",
			"  todolist update [global options] <todo> [<title>] [<status>] [<priority>]",
			"  todolist update [global options] <todo> [title=<title>] [status=<status>] [priority=<priority>]",
			"",
		}, "\n")
	case "delete":
		return strings.Join([]string{
			"Usage:",
			"  todolist delete [global options] <todo>",
			"",
		}, "\n")
	default:
		return rootHelp()
	}
}
