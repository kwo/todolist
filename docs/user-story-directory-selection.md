# User Story: task directory selection

Add directory selection controls to tasklist.

## User story

As a user,
I want to choose which task directory the CLI operates on,
so that I can work with multiple task sets.

## Goal

This work adds a global directory option and environment variable support.

## Command surface

Global option:

```bash
tasklist <command> -d <dir>
tasklist <command> --directory <dir>
```

Environment variable:

```bash
TASKLIST_DIRECTORY=<dir> tasklist <command>
```

The CLI is command-first, so `--directory` appears after the command.

Example:

```bash
tasklist list -d ./work-tasks
```

## Acceptance criteria

Given:

- tasklist needs to resolve a task directory

Then:

- resolution order is:
  1. `-d, --directory`
  2. `TASKLIST_DIRECTORY`
  3. default `./tasks`

## Scope

- global `-d, --directory` option
- `TASKLIST_DIRECTORY` environment variable
- directory resolution order
- directory-related error handling

## Open issues

1. The interaction between directory selection and initialization workflows should be documented alongside `init`.
