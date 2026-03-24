## Using the todolist app

Use the `todolist` CLI to manage todos. Prefer the CLI over editing todo files directly.

### Directory resolution

The todo directory is chosen in this order:

1. `-d` / `--directory`
2. `TODOLIST_DIRECTORY`
3. `./todo`

Important: global flags can appear anywhere after the command name.

Examples:

```bash
todolist list -d ./todo
todolist add -d ./work-todos --title "Buy groceries"
todolist list --json
```

### Initialize a directory

Before using a new todo directory, run:

```bash
todolist init
```

This creates the todo directory and a `.todos` config file if they do not already exist.

### Todo fields

Each todo has:

- `id`
- `title`
- `status`
- `priority`
- `createdAt`
- `lastModified`
- `description`

Valid statuses:

- `todo`
- `wip`
- `done`

Valid priorities:

- `1` to `5`, where `1` is highest and `5` is lowest

Defaults:

- status = `todo`
- priority = `5`

### Add a todo

```bash
todolist add -t <title> [-s <status>] [-p <priority>]
todolist add --title <title> [--status <status>] [--priority <priority>]
```

Examples:

```bash
todolist add --title "Buy groceries"
todolist add -t "Buy groceries" -s wip -p 2
todolist add --title "Buy groceries" --status wip --priority 2
```

A single positional argument is also accepted as the title:

```bash
todolist add "Buy groceries"
```

Use `-t` / `--title` when the title could be mistaken for another value:

```bash
todolist add -t "done"
todolist add --title "2"
```

Non-JSON output prints only the new todo ID.

### Add or replace description via stdin

Pipe Markdown on stdin:

```bash
printf 'Need milk, eggs, and bread.\n' | todolist add --title "Buy groceries" --status wip --priority 2
printf 'Need milk, eggs, bread, and chips.\n' | todolist update todo-7k9m --title "Updated title"
```

Rules:

- `add`: stdin becomes the description
- `update`: stdin replaces the full description
- if stdin is not provided to `update`, the description stays unchanged

### List todos

```bash
todolist list [-s <status-filter>] [-p <priority-filter>]
todolist list [--status <status-filter>] [--priority <priority-filter>]
```

Default behavior:

- `todolist list` excludes `done` todos

Examples:

```bash
todolist list
todolist list -s done
todolist list --status done
todolist list --status done!
todolist list -p 1
todolist list --priority 1
todolist list --priority 3-
todolist list --priority 3+
todolist list --priority 3!
todolist list -s done -p 3+
todolist list --status done --priority 3+
```

Filter meanings:

- `done` = only `done`
- `done!` = exclude `done`
- `1` = only priority 1
- `3-` = priorities less than 3
- `3+` = priorities greater than 3
- `3!` = priorities not equal to 3

Text output columns:

```text
<id>\t<priority>\t<status>\t<title>
```

Text list output truncates long titles. Use `view` or `--json` when full values matter.

### View a todo

```bash
todolist view <todo-id>
todolist view --json <todo-id>
```

- text output returns the raw Markdown file
- JSON output returns the parsed todo object

### Update a todo

```bash
todolist update <todo-id> [-t <title>] [-s <status>] [-p <priority>]
todolist update <todo-id> [--title <title>] [--status <status>] [--priority <priority>]
```

Examples:

```bash
todolist update todo-7k9m -t "Buy groceries and snacks"
todolist update todo-7k9m --title "Buy groceries and snacks"
todolist update todo-7k9m -s done -p 1
todolist update todo-7k9m --status done --priority 1
todolist update todo-7k9m --title "Buy groceries" --status wip --priority 2
```

Rule:

- `update` must change at least one of title, status, priority, or description via stdin

### Delete a todo

```bash
todolist delete <todo-id>
```

### Prefer JSON for agent automation

Use `--json` whenever an agent needs structured output:

```bash
todolist list --json
todolist view --json todo-7k9m
todolist add --json --title "Buy groceries"
todolist update --json todo-7k9m --status done
todolist delete --json todo-7k9m
```

Recommended agent workflow:

1. `todolist list --json`
2. `todolist view --json <id>`
3. `todolist update --json <id> --status done` or `todolist delete --json <id>`

### File format

Todos are stored as `<id>.md` Markdown files with YAML front matter and a Markdown description body.

Prefer the CLI for normal changes. If manual editing is necessary:

- keep valid YAML front matter
- keep timestamps in RFC 3339 UTC
- keep status within `todo|wip|done`
- keep priority within `1..5`
