---
id: todo-2yh4
title: 'User Story: machine-readable JSON output'
status: done
createdAt: "2026-03-23T18:50:57Z"
lastModified: "2026-03-23T18:50:57Z"
---

# User Story: machine-readable JSON output

Add machine-readable JSON output to todolist.

## User story

As a user,
I want command output in JSON,
so that I can script against todolist and integrate it with other tools.

## Goal

This work adds a global `--json` option for supported commands.

## Command surface

Global option:

```bash
todolist <command> --json
```

Examples:

```bash
todolist list --json
todolist view --json todo-7k9m
todolist add --json "Buy groceries"
```

## Acceptance criteria

Given:

- the user passes `--json`

Then:

- command results are written as JSON instead of human-readable text
- the option is global and, in the command-first CLI, appears after the command

Supported commands:

- `add`
- `list`
- `view`
- `update`
- `delete`

## Scope

- global `--json` option
- JSON output for core commands

## JSON shape

Todo objects are encoded as JSON objects with these fields. The `description` field is omitted when it is empty:

```json
{
  "id": "todo-7k9m",
  "title": "Buy groceries",
  "status": "todo",
  "priority": 2,
  "createdAt": "2026-03-18T10:00:00Z",
  "lastModified": "2026-03-18T10:00:00Z",
  "description": "Need milk, eggs, and bread.\n"
}
```

Command outputs:

- `add` → the created todo object
- `list` → an array of todo objects
- `view` → the requested todo object
- `update` → the updated todo object
- `delete` → an object of the form:

```json
{
  "id": "todo-7k9m",
  "deleted": true
}
```

## Open issues

None currently.
