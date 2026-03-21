# User Story: machine-readable JSON output

Add machine-readable JSON output to tasklist.

## User story

As a user,
I want command output in JSON,
so that I can script against tasklist and integrate it with other tools.

## Goal

This work adds a global `--json` option for supported commands.

## Command surface

Global option:

```bash
tasklist <command> --json
```

Examples:

```bash
tasklist list --json
tasklist view --json task-7k9m
tasklist add --json "Buy groceries"
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

## Open issues

1. The exact JSON shape for each command should be specified before implementation.
