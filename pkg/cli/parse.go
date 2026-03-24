package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kwo/todolist/pkg/todolist"
)

//nolint:maintidx // Command dispatch is intentionally centralized in one parser.
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
	case "init":
		if parsed.help {
			return parsed, nil
		}

		runner, initErr := parseNoValueCommand(commandValues)
		if initErr != nil {
			return parsedCommand{}, initErr
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
	case "usage":
		if parsed.help {
			return parsed, nil
		}

		runner, usageErr := parseUsageCommand(commandValues)
		if usageErr != nil {
			return parsedCommand{}, usageErr
		}

		parsed.runner = runner
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

func parseNoValueCommand(values []string) (initCommand, error) {
	if len(values) == 0 {
		return initCommand{}, nil
	}

	return initCommand{}, fmt.Errorf("cannot assign value %q", values[0])
}

func parseUsageCommand(values []string) (usageCommand, error) {
	if len(values) == 0 {
		return usageCommand{}, nil
	}

	return usageCommand{}, fmt.Errorf("cannot assign value %q", values[0])
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
