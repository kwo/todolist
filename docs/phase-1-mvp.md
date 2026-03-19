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

Behavior:

- show a compact human-readable list
- for now it shows the task id and title

### `view`

Show a single task.

```bash
tasklist view <task>
```

Arguments:

- `<task>` — required task ID

Flags:

- none in the MVP

Behavior:

- show all task metadata and the raw Markdown description
- do not render Markdown

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
- a successful update should automatically set `lastModified` to the current time unless explicitly overridden in a later phase

Examples:

```bash
tasklist update task-7k9m --title "Buy groceries and snacks"
printf 'Need milk, eggs, bread, and chips.\n' | tasklist update task-7k9m
printf 'Need milk, eggs, bread, and chips.\n' | tasklist update task-7k9m --title "Buy groceries and snacks"
```

### `delete`

Delete a task.

```bash
tasklist delete <task> [-f]
```

Arguments:

- `<task>` — required task ID

Flags:

- `-f`, `--force` — skip confirmation if the CLI prompts before deletion

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
- `createdAt` — set automatically by the CLI by default
- `lastModified` — set automatically by the CLI by default and updated automatically on every successful `update`

The Markdown description is also part of the task and is updatable.

## Scope

- Markdown task files with YAML front matter
- one task per file
- basic CRUD commands
- auto-generated task IDs
- file naming as `<id>.md`
- default task directory support
- human-readable CLI output
- description input through stdin for `add` and `update`

Shared specifications such as ID generation, file naming, principles, and directory selection are documented in [README.md](../README.md).

## Open Issues

None currently.
