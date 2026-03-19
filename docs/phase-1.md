# Phase 1: Core task management polish

Phase 1 builds on the MVP by adding richer task metadata, filtering, machine-readable output, and directory/config controls.

It should preserve the MVP command flow:

- use **arguments** for the main thing a command acts on
- use **flags** for optional metadata and behavior changes
- continue to read the Markdown description from **stdin** when stdin is piped

## Commands

Global options in Phase 1:

- `--json` — output machine-readable JSON
- `-d <dir>`, `--directory <dir>` — use a specific task directory

Because `--directory` is global, it should appear before the subcommand. Example:

```bash
tasklist -d ./work-tasks list
```

### `add`

Create a task.

```bash
tasklist [--json] [-d <dir>] add <title> [--status <status>] [--priority <priority>]
```

Arguments:

- `<title>` — required task title

Flags:

- `--status <status>` — set the task status
- `--priority <priority>` — set the task priority

Description input:

- if stdin is piped, the CLI should read the task description from stdin
- if stdin is not piped, the task should be created with an empty description

### `list`

List tasks.

```bash
tasklist [--json] [-d <dir>] list [--priority <priority>]
```

Arguments:

- none

Flags:

- `--priority <priority>` — filter tasks by priority

Behavior:

- retain the MVP’s compact human-readable list by default
- for now it shows the task id and title

### `view`

Show a single task.

```bash
tasklist [--json] [-d <dir>] view <task>
```

Arguments:

- `<task>` — required task ID

Flags:

- none command-specific in Phase 1

Behavior:

- show all task metadata and the raw Markdown description
- do not render Markdown

### `update`

Update a task.

```bash
tasklist [--json] [-d <dir>] update <task> [--title <title>] [--status <status>] [--priority <priority>]
```

Arguments:

- `<task>` — required task ID

Flags:

- `--title <title>` — replace the task title
- `--status <status>` — replace the task status
- `--priority <priority>` — replace the task priority

Description input:

- if stdin is piped, the CLI should replace the task description with stdin contents
- if stdin is not piped, the description should remain unchanged

Behavior:

- `id` is not updatable
- the command should require at least one change, either a field flag or piped stdin

### `delete`

Delete a task.

```bash
tasklist [--json] [-d <dir>] delete <task> [-f]
```

Arguments:

- `<task>` — required task ID

Flags:

- `-f`, `--force` — skip confirmation if the CLI prompts before deletion

## Task Format Additions

Phase 1 extends the task front matter with:

- `status`
- `priority`

Example:

```md
---
id: task-7k9m
title: Buy groceries
status: todo
priority: 2
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

## Scope

- `status` metadata
- `priority` metadata
- filtering in `list` by priority
- JSON output via global `--json`
- directory selection via `-d, --directory` or `TASKLIST_DIRECTORY`
- configurable task ID prefixes via `.tasks`
- continued stdin-based description input for `add` and `update`

Shared specifications such as status values, priority values, ID generation, file naming, principles, and directory selection are documented in [README.md](../README.md).

## Open Issues

1. `update` should auto-update `lastModified`. It's not stated that `lastModified` should be automatically updated when a task is changed via `update`.
