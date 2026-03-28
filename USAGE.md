## Using the todolist app

Use the `todolist` CLI to manage todos. Prefer the CLI over editing todo files directly.

### Directory resolution

The todo directory is always:

1. `./todo`

Examples:

```bash
todolist list
todolist add --title "Buy groceries"
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
- `parents`
- `depends`
- `createdAt`
- `lastModified`
- `description`

Computed output fields:

- `ready` = `true` when all dependencies are `done`; otherwise `false`

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
todolist add -t <title> [-s <status>] [-p <priority>] [--parent <todo-id> ...] [--depends <todo-id> ...]
todolist add --title <title> [--status <status>] [--priority <priority>] [--parent <todo-id> ...] [--depends <todo-id> ...]
```

Examples:

```bash
todolist add --title "Buy groceries"
todolist add -t "Buy groceries" -s wip -p 2
todolist add --title "Buy groceries" --status wip --priority 2
todolist add --title "Buy groceries" --parent todo-3h7q --parent todo-9p2d
todolist add --title "Buy groceries" --depends todo-3h7q --depends todo-9p2d
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
todolist list [-s <status-filter>] [-p <priority-filter>] [--all]
todolist list [--status <status-filter>] [--priority <priority-filter>] [--all]
```

Default behavior:

- `todolist list` excludes `done` todos
- `todolist list` includes only todos whose computed `ready` is `true`
- list results are sorted by priority ascending, then title ascending
- `--all` disables readiness filtering so both ready and blocked todos are shown

Examples:

```bash
todolist list
todolist list --all
todolist list -s done
todolist list --status done
todolist list --status done!
todolist list -p 1
todolist list --priority 1
todolist list --priority 3-
todolist list --priority 3+
todolist list --priority 3!
todolist list -s done -p 3+ --all
todolist list --status done --priority 3+
todolist list --status wip --priority 2 --all
```

Filter meanings:

- `done` = only `done`
- `done!` = exclude `done`
- `1` = only priority 1
- `3-` = priorities less than 3
- `3+` = priorities greater than 3
- `3!` = priorities not equal to 3
- `--all` = include both ready and blocked todos

Text output columns:

```text
<id>\t<priority>\t<status>\t<title>\t<first-parent-id>\t<first-dependency-id>
```

If a todo has multiple parents, the parent column shows the first parent ID followed by `,...`.

If a todo has multiple dependencies, the dependency column shows the first dependency ID followed by `,...`.

Text list output truncates long titles. Use `view` or `--json` when full values matter.

JSON list output continues to include computed `ready`.

### View a todo

```bash
todolist view <todo-id>
todolist view --json <todo-id>
```

- text output returns the todo Markdown plus a human-friendly `Parents:` section when parents exist
- each parent in the human-friendly section is shown on a single line as `- <id> <title>`
- JSON output returns the parsed todo object including computed `ready`

### Show version information

```bash
todolist version
todolist version --json
```

Text output prints only the version string.

JSON output returns:

- `version`
- `commit` when embedded VCS metadata is available, shortened to 8 characters
- `dirty` when embedded VCS metadata reports a dirty working tree
- `runtime` when Go build metadata is available

### Update a todo

```bash
todolist update <todo-id> [-t <title>] [-s <status>] [-p <priority>] [--parent <todo-id>|<todo-id>! ...] [--depends <todo-id>|<todo-id>! ...]
todolist update <todo-id> [--title <title>] [--status <status>] [--priority <priority>] [--parent <todo-id>|<todo-id>! ...] [--depends <todo-id>|<todo-id>! ...]
```

Examples:

```bash
todolist update todo-7k9m -t "Buy groceries and snacks"
todolist update todo-7k9m --title "Buy groceries and snacks"
todolist update todo-7k9m -s done -p 1
todolist update todo-7k9m --status done --priority 1
todolist update todo-7k9m --title "Buy groceries" --status wip --priority 2
todolist update todo-7k9m --parent todo-3h7q --parent todo-9p2d
todolist update todo-7k9m --parent todo-3h7q!
todolist update todo-7k9m --depends todo-3h7q --depends todo-9p2d
todolist update todo-7k9m --depends todo-3h7q!
```

Rules:

- `update` must change at least one of title, status, priority, parents, dependencies, or description via stdin
- `--parent <todo-id>` adds a parent
- `--parent <todo-id>!` removes a parent
- removing a parent that is not currently assigned fails
- `--depends <todo-id>` adds a dependency
- `--depends <todo-id>!` removes a dependency
- duplicate dependency additions are deduplicated
- removing a dependency that is not currently assigned fails

### Delete a todo

```bash
todolist delete <todo-id>
```

If other todos list the deleted todo in `parents` or `depends`, the delete still succeeds and those references are removed from those child todos.

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
- keep `parents` as a YAML list of existing todo IDs
- keep `depends` as a YAML list of existing todo IDs
- do not store `ready`; it is computed at read time
