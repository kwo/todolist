# Phase 1 MVP: Basic task management

The MVP focuses on the smallest useful version of Task List: basic task CRUD backed by Markdown files.

## CLI shape

For the MVP, the CLI should optimize for fast human typing:

- use **arguments** for the main thing a command acts on
- use **flags** only for optional behavior changes
- use **stdin** for Markdown description input instead of a dedicated description flag

That means:

- task IDs should be positional arguments for `view`, `update`, and `delete`
- the task title should be a positional argument for `add`
- `update` should use flags only for fields that are optional to change
- the task description should be read from stdin when stdin is piped

## Commands

### `add`

Create a task.

```bash
tasklist add <title>
```

Arguments:

- `<title>` — required task title

Flags:

- none in the MVP

Description input:

- if stdin is piped, the CLI should read the task description from stdin
- if stdin is not piped, the task should be created with an empty description

Output:

- on success, print the generated task ID to stdout
- on error, print a message to stderr and exit with code 1

Examples:

```bash
tasklist add "Buy groceries"
printf 'Need milk, eggs, and bread.\n' | tasklist add "Buy groceries"
```

### `list`

List tasks.

```bash
tasklist list
```

Arguments:

- none

Flags:

- none in the MVP

Output:

- one task per line, tab-separated: `<id>\t<title>`
- if there are no tasks, print nothing

Example output:

```
task-7k9m	Buy groceries
task-2w8x	Write proposal
```

### `view`

Show a single task.

```bash
tasklist view <task>
```

Arguments:

- `<task>` — required task ID

Flags:

- none in the MVP

Output:

- print the raw task file as-is, including the YAML front matter and Markdown description
- do not render Markdown
- if the task does not exist, print an error to stderr and exit with code 1

Example output:

```
---
id: task-7k9m
title: Buy groceries
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

### `update`

Update a task.

```bash
tasklist update <task> [--title <title>]
```

Arguments:

- `<task>` — required task ID

Flags:

- `--title <title>` — replace the task title

Description input:

- if stdin is piped, the CLI should replace the task description with stdin contents
- if stdin is not piped, the description should remain unchanged

Behavior:

- `id` is not updatable
- the command should require at least one change, either `--title` or piped stdin
- a successful update should automatically set `lastModified` to the current time

Output:

- on success, print nothing
- if the task does not exist, print an error to stderr and exit with code 1
- if no changes are provided, print an error to stderr and exit with code 1

Examples:

```bash
tasklist update task-7k9m --title "Buy groceries and snacks"
printf 'Need milk, eggs, bread, and chips.\n' | tasklist update task-7k9m
printf 'Need milk, eggs, bread, and chips.\n' | tasklist update task-7k9m --title "Buy groceries and snacks"
```

### `delete`

Delete a task.

```bash
tasklist delete <task>
```

Arguments:

- `<task>` — required task ID

Flags:

- none in the MVP

Behavior:

- delete the task file immediately without prompting
- confirmation prompting and `--force` are deferred to Phase 2, where parent-child relationships introduce cases that need protection

Output:

- on success, print nothing
- if the task does not exist, print an error to stderr and exit with code 1

## Task Format

Each task is a Markdown file with YAML front matter and an optional Markdown description.

Example:

```md
---
id: task-7k9m
title: Buy groceries
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

Front matter fields in the MVP:

- `id` — unique task identifier; not updatable
- `title`
- `createdAt` — set on creation to the current time; uses RFC 3339 in UTC
- `lastModified` — set on creation to the current time; updated automatically on every successful `update`; uses RFC 3339 in UTC

On creation, `createdAt` and `lastModified` should both be set to the same current time.

The Markdown description is also part of the task and is updatable.

The MVP does not write `status` or `priority` fields to task files. Those are introduced in Phase 1.

## Error Handling

- nonexistent task ID → error message to stderr, exit code 1
- missing or inaccessible task directory → error message to stderr, exit code 1
- `update` with no changes provided → error message to stderr, exit code 1
- all other errors → message to stderr, exit code 1
- success → exit code 0

## Deferred to later phases

The following are explicitly **not** part of the MVP:

- `status` and `priority` metadata (Phase 1)
- filtering in `list` (Phase 1)
- `--json` global option (Phase 1)
- `-d, --directory` global option (Phase 1)
- `TASKLIST_DIRECTORY` environment variable (Phase 1)
- `.tasks` config file and configurable ID prefix (Phase 1)
- confirmation prompting and `--force` on `delete` (Phase 2)
- parent-child relationships (Phase 2)
- dependencies (Phase 2)

In the MVP, the task directory is always `./tasks` and the ID prefix is always `task-`.

## Scope

- Markdown task files with YAML front matter
- one task per file
- basic CRUD commands: `add`, `list`, `view`, `update`, `delete`
- auto-generated task IDs with hardcoded `task-` prefix
- file naming as `<id>.md`
- task directory is `./tasks`
- human-readable CLI output
- description input through stdin for `add` and `update`
- timestamps in RFC 3339 UTC
- exit code 0 on success, 1 on error

Shared specifications such as ID generation, file naming, and principles are documented in [README.md](../README.md).

## Open Issues

None currently.
