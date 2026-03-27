package cli

import (
	"fmt"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/kwo/todolist/pkg/todolist"
)

func parseArgs(args []string, app *App) (parsedCommand, error) {
	if len(args) == 0 {
		return parsedCommand{}, fmt.Errorf("missing command")
	}

	if isHelpToken(args[0]) {
		return parsedCommand{name: "", help: true, globals: resolveRunOptions(app, globalOptions{})}, nil
	}

	commandName := args[0]
	commandArgs := args[1:]

	switch commandName {
	case "add":
		return parseAddArgs(commandArgs, app)
	case "init":
		return parseNoArgCommand(commandName, commandArgs, app, initCommand{})
	case "list":
		return parseListArgs(commandArgs, app)
	case "view":
		return parseSingleIDCommand(commandName, commandArgs, app, func(id string) commandRunner { return viewCommand{Todo: id} })
	case "update":
		return parseUpdateArgs(commandArgs, app)
	case "delete":
		return parseSingleIDCommand(commandName, commandArgs, app, func(id string) commandRunner { return deleteCommand{Todo: id} })
	case "usage":
		return parseNoArgCommand(commandName, commandArgs, app, usageCommand{})
	default:
		return parsedCommand{}, fmt.Errorf("unknown command %q", commandName)
	}
}

// newFlagSet creates a new pflag.FlagSet for a command with global flags registered.
func newFlagSet(name string, globals *globalOptions) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.StringVarP(&globals.Directory, "directory", "d", "", "use a specific todo directory")
	fs.BoolVar(&globals.JSON, "json", false, "enable JSON output")
	fs.BoolVarP(&globals.Help, "help", "h", false, "show help")

	return fs
}

func parseNoArgCommand(name string, args []string, app *App, runner commandRunner) (parsedCommand, error) {
	var globals globalOptions

	fs := newFlagSet(name, &globals)
	if err := fs.Parse(args); err != nil {
		return parsedCommand{}, err
	}

	parsed := parsedCommand{
		name:    name,
		help:    globals.Help,
		globals: resolveRunOptions(app, globals),
	}

	if parsed.help {
		return parsed, nil
	}

	if fs.NArg() > 0 {
		return parsedCommand{}, fmt.Errorf("%s does not accept positional arguments, got %q", name, fs.Arg(0))
	}

	parsed.runner = runner

	return parsed, nil
}

func parseSingleIDCommand(name string, args []string, app *App, makeRunner func(string) commandRunner) (parsedCommand, error) {
	var globals globalOptions

	fs := newFlagSet(name, &globals)
	if err := fs.Parse(args); err != nil {
		return parsedCommand{}, err
	}

	parsed := parsedCommand{
		name:    name,
		help:    globals.Help,
		globals: resolveRunOptions(app, globals),
	}

	if parsed.help {
		return parsed, nil
	}

	if fs.NArg() == 0 {
		return parsedCommand{}, fmt.Errorf("%s requires a todo id", name)
	}

	if fs.NArg() > 1 {
		return parsedCommand{}, fmt.Errorf("%s accepts only one positional argument, got extra %q", name, fs.Arg(1))
	}

	parsed.runner = makeRunner(fs.Arg(0))

	return parsed, nil
}

func parseAddArgs(args []string, app *App) (parsedCommand, error) {
	var globals globalOptions

	var title string

	var status string

	var priority int
	var parents []string
	var depends []string

	fs := newFlagSet("add", &globals)
	fs.StringVarP(&title, "title", "t", "", "todo title (required)")
	fs.StringVarP(&status, "status", "s", todolist.DefaultStatus, "todo status: todo|wip|done")
	fs.IntVarP(&priority, "priority", "p", todolist.DefaultPriority, "priority 1..5")
	fs.StringArrayVar(&parents, "parent", nil, "assign a parent todo ID; repeat to add multiple parents")
	fs.StringArrayVar(&depends, "depends", nil, "assign a dependency todo ID; repeat to add multiple dependencies")

	if err := fs.Parse(args); err != nil {
		return parsedCommand{}, err
	}

	parsed := parsedCommand{
		name:    "add",
		help:    globals.Help,
		globals: resolveRunOptions(app, globals),
	}

	if parsed.help {
		return parsed, nil
	}

	// Allow a single positional arg as the title when --title is not set.
	if title == "" {
		if fs.NArg() == 0 {
			return parsedCommand{}, fmt.Errorf("add requires --title or a positional title argument")
		}

		if fs.NArg() > 1 {
			return parsedCommand{}, fmt.Errorf("add accepts only one positional argument (title), got extra %q; use --status and --priority flags", fs.Arg(1))
		}

		title = fs.Arg(0)
	} else if fs.NArg() > 0 {
		return parsedCommand{}, fmt.Errorf("cannot use both --title flag and positional title argument %q", fs.Arg(0))
	}

	if err := todolist.ValidateStatus(status); err != nil {
		return parsedCommand{}, err
	}

	if err := todolist.ValidatePriority(priority); err != nil {
		return parsedCommand{}, err
	}

	parsed.runner = addCommand{
		Title:    title,
		Status:   status,
		Priority: priority,
		Parents:  parents,
		Depends:  depends,
	}

	return parsed, nil
}

func parseUpdateArgs(args []string, app *App) (parsedCommand, error) {
	var globals globalOptions

	var title string

	var status string

	var priority int
	var parents []string
	var depends []string

	fs := newFlagSet("update", &globals)
	fs.StringVarP(&title, "title", "t", "", "new title")
	fs.StringVarP(&status, "status", "s", "", "new status: todo|wip|done")
	fs.IntVarP(&priority, "priority", "p", 0, "new priority 1..5")
	fs.StringArrayVar(&parents, "parent", nil, "add a parent todo ID, or remove one with a trailing !")
	fs.StringArrayVar(&depends, "depends", nil, "add a dependency todo ID, or remove one with a trailing !")

	if err := fs.Parse(args); err != nil {
		return parsedCommand{}, err
	}

	parsed := parsedCommand{
		name:    "update",
		help:    globals.Help,
		globals: resolveRunOptions(app, globals),
	}

	if parsed.help {
		return parsed, nil
	}

	if fs.NArg() == 0 {
		return parsedCommand{}, fmt.Errorf("update requires a todo id")
	}

	if fs.NArg() > 1 {
		return parsedCommand{}, fmt.Errorf("update accepts only one positional argument (todo id), got extra %q; use --title, --status, --priority flags", fs.Arg(1))
	}

	cmd := updateCommand{Todo: fs.Arg(0)}

	if fs.Changed("title") {
		cmd.Title = title
		cmd.TitleProvided = true
	}

	if fs.Changed("status") {
		if err := todolist.ValidateStatus(status); err != nil {
			return parsedCommand{}, err
		}

		cmd.Status = status
		cmd.StatusProvided = true
	}

	if fs.Changed("priority") {
		if err := todolist.ValidatePriority(priority); err != nil {
			return parsedCommand{}, err
		}

		cmd.Priority = priority
		cmd.PriorityProvided = true
	}

	if fs.Changed("parent") {
		cmd.Parents = parents
		cmd.ParentsProvided = true
	}

	if fs.Changed("depends") {
		cmd.Depends = depends
		cmd.DependsProvided = true
	}

	parsed.runner = cmd

	return parsed, nil
}

func parseListArgs(args []string, app *App) (parsedCommand, error) {
	var globals globalOptions

	var status string

	var priority string

	fs := newFlagSet("list", &globals)
	fs.StringVarP(&status, "status", "s", "", "status filter: todo|wip|done, append ! to exclude")
	fs.StringVarP(&priority, "priority", "p", "", "priority filter: n, n!, n+, or n-")

	if err := fs.Parse(args); err != nil {
		return parsedCommand{}, err
	}

	parsed := parsedCommand{
		name:    "list",
		help:    globals.Help,
		globals: resolveRunOptions(app, globals),
	}

	if parsed.help {
		return parsed, nil
	}

	if fs.NArg() > 0 {
		return parsedCommand{}, fmt.Errorf("list does not accept positional arguments, got %q; use --status and --priority flags", fs.Arg(0))
	}

	cmd := listCommand{}

	if fs.Changed("status") {
		statusFilter, excludeStatus, err := parseStatusFilter(status)
		if err != nil {
			return parsedCommand{}, err
		}

		cmd.StatusFilter = statusFilter
		cmd.ExcludeStatus = excludeStatus
	} else {
		// Default: exclude done.
		cmd.StatusFilter = "done"
		cmd.ExcludeStatus = true
	}

	if fs.Changed("priority") {
		if _, err := parsePriorityFilter(priority); err != nil {
			return parsedCommand{}, err
		}

		cmd.PriorityFilter = strings.TrimSpace(priority)
	}

	parsed.runner = cmd

	return parsed, nil
}

func parseStatusFilter(raw string) (string, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, fmt.Errorf("invalid status %q: must be one of todo, wip, done", value)
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

func isHelpToken(value string) bool {
	return value == "-h" || value == "--help"
}
