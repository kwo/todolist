---
id: todo-mes9
title: 'User Story: todo directory selection'
status: done
createdAt: "2026-03-23T18:50:56Z"
lastModified: "2026-03-23T18:50:56Z"
---

# User Story: todo directory selection

Add directory selection controls to todolist.

## User story

As a user,
I want to choose which todo directory the CLI operates on,
so that I can work with multiple todo sets.

## Goal

This work adds a global directory option and environment variable support.

## Command surface

Global option:

```bash
todolist <command> -d <dir>
todolist <command> --directory <dir>
```

Environment variable:

```bash
TODOLIST_DIRECTORY=<dir> todolist <command>
```

The CLI is command-first, so `--directory` appears after the command.

Example:

```bash
todolist list -d ./work-todos
```

## Acceptance criteria

Given:

- todolist needs to resolve a todo directory

Then:

- resolution order is:
  1. `-d, --directory`
  2. `TODOLIST_DIRECTORY`
  3. default `./todo`

## Scope

- global `-d, --directory` option
- `TODOLIST_DIRECTORY` environment variable
- directory resolution order
- directory-related error handling

## Open issues

1. The interaction between directory selection and initialization workflows should be documented alongside `init`.
