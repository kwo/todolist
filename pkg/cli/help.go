package cli

import (
	"fmt"
	"io"
	"strings"
)

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
		"  init",
		"  list [<status-filter>] [<priority-filter>]",
		"  view <todo>",
		"  update <todo> [<title>] [<status>] [<priority>]",
		"  delete <todo>",
		"  usage",
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
	case "init":
		return strings.Join([]string{
			"Usage:",
			"  todolist init [global options]",
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
	case "usage":
		return strings.Join([]string{
			"Usage:",
			"  todolist usage [global options]",
			"",
		}, "\n")
	default:
		return rootHelp()
	}
}
