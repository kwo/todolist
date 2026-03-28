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
		"  todolist <command> [flags]",
		"",
		"Commands:",
		"  add      add a new todo",
		"  init     initialize a todo directory",
		"  list     list todos",
		"  view     view a todo",
		"  update   update a todo",
		"  delete   delete a todo",
		"  usage    print usage documentation",
		"  version  print version information",
		"",
		"Global flags:",
		"      --json  enable JSON output",
		"  -h, --help  show help",
		"",
	}, "\n")
}

func commandHelp(command string) string {
	switch command {
	case "add":
		return strings.Join([]string{
			"Usage:",
			"  todolist add [flags] [<title>]",
			"",
			"Flags:",
			"  -t, --title <title>     todo title (required if no positional title)",
			"  -s, --status <status>   todo status: todo|wip|done (default \"todo\")",
			"  -p, --priority <n>      priority 1..5 (default 5)",
			"      --depends <id>      dependency todo ID; repeat to add multiple",
			"",
			"A single positional argument is accepted as the title.",
			"Use --title when the title could be mistaken for another value.",
			"",
		}, "\n")
	case "init":
		return strings.Join([]string{
			"Usage:",
			"  todolist init [flags]",
			"",
		}, "\n")
	case "list":
		return strings.Join([]string{
			"Usage:",
			"  todolist list [flags]",
			"",
			"Flags:",
			"  -s, --status <filter>     status filter: todo|wip|done, append ! to exclude",
			"  -p, --priority <filter>   priority filter: n, n!, n+, or n-",
			"      --all                include both ready and blocked todos",
			"",
			"By default, done todos are excluded and only ready todos are shown.",
			"",
		}, "\n")
	case "view":
		return strings.Join([]string{
			"Usage:",
			"  todolist view [flags] <todo-id>",
			"",
		}, "\n")
	case "update":
		return strings.Join([]string{
			"Usage:",
			"  todolist update [flags] <todo-id>",
			"",
			"Flags:",
			"  -t, --title <title>     new title",
			"  -s, --status <status>   new status: todo|wip|done",
			"  -p, --priority <n>      new priority 1..5",
			"      --depends <id>      add a dependency todo ID, or remove one with !",
			"",
			"At least one of --title, --status, --priority, --parent, --depends, or stdin description is required.",
			"",
		}, "\n")
	case "delete":
		return strings.Join([]string{
			"Usage:",
			"  todolist delete [flags] <todo-id>",
			"",
		}, "\n")
	case "usage":
		return strings.Join([]string{
			"Usage:",
			"  todolist usage [flags]",
			"",
		}, "\n")
	case "version":
		return strings.Join([]string{
			"Usage:",
			"  todolist version [flags]",
			"",
		}, "\n")
	default:
		return rootHelp()
	}
}
