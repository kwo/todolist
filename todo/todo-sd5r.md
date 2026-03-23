---
id: todo-sd5r
title: 'User Story: list defaults to done! status filter'
status: todo
priority: 3
createdAt: "2026-03-23T18:59:19Z"
lastModified: "2026-03-23T18:59:19Z"
---

# User Story: list defaults to `done!` status filter

Change the default behavior of the `list` command so that completed todos are hidden unless the user explicitly provides a status filter.

## User story

As a user,
I want `todolist list` to exclude completed todos by default,
so that I can focus on active work without manually adding a filter each time.

## Goal

This work changes the default status filtering behavior of `list`:

- if the user does not provide any status filter argument, `list` behaves as if `done!` had been supplied
- if the user explicitly provides a status filter, that explicit filter takes precedence

## Command surface

### `list`

```bash
todolist list
todolist list done
todolist list done!
todolist list status=done
todolist list 2
todolist list priority=3+
```

## Behavior

### Default case

When the user runs:

```bash
todolist list
```

Then the command should behave like:

```bash
todolist list 'done!'
```

So todos with `status: done` are excluded by default.

### Explicit status filters

When the user provides an explicit status filter, the explicit filter should be used instead of the default.

Examples:

```bash
todolist list done
todolist list 'done!'
todolist list status=done
todolist list status='done!'
```

### Priority-only filters

If the user provides only a priority filter and no status filter, the default `done!` filter still applies.

Examples:

```bash
todolist list 2
todolist list 3+
todolist list priority=3-
```

These should behave as though the user had also provided `done!`.

## Acceptance criteria

### Default hidden completed todos

Given:

- the todo store contains todos with statuses `todo`, `wip`, and `done`

When:

- the user runs `todolist list`

Then:

- todos with `status: done` are excluded from the output
- todos with statuses other than `done` remain visible

### Explicit status overrides default

Given:

- the user provides a status filter to `list`

When:

- the command runs

Then:

- the explicit status filter is used
- the implicit `done!` default is not applied in addition to it

### Priority-only input still uses default status filter

Given:

- the user provides only a priority filter to `list`

When:

- the command runs

Then:

- the priority filter is applied
- the implicit `done!` status filter is also applied

### JSON compatibility

Given:

- the user runs `todolist list --json`

Then:

- the same default filtering behavior applies
- only the output format differs

## Scope

- default status filtering behavior for `list`
- precedence rules between implicit default filtering and explicit status filters
- compatibility with existing priority filters and JSON output

## Open issues

1. It should be decided whether a future flag or config option should allow users to opt out of this default behavior.
