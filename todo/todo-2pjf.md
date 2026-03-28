---
id: todo-2pjf
title: replace list --ready filter with valueless --all flag
status: todo
priority: 5
createdAt: "2026-03-28T19:52:09Z"
lastModified: "2026-03-28T19:52:09Z"
---

# User Story: replace `--ready` list filtering with `--all`

Change `todolist list` so the `--ready` flag is replaced by a valueless `--all` flag whose presence means list all matching todos regardless of readiness.

## User story

As a user,
I want a simple `--all` flag on `todolist list`,
so that I can easily include blocked and ready todos in one listing without passing a boolean value.

## Goal

This story simplifies the list command surface by removing the boolean `--ready <true|false>` flag and replacing it with a presence-based `--all` flag. By default, `todolist list` should continue to show only ready, non-done todos. When `--all` is present, readiness filtering should be disabled and both ready and blocked todos should be shown.

## Acceptance criteria

### Default behavior remains ready-only

Given:

- the user runs `todolist list` with no `--all` flag

Then:

- the command excludes `done` todos by default
- the command includes only todos whose computed `ready` value is `true`

### `--all` includes ready and blocked todos

Given:

- the user runs `todolist list --all`

Then:

- the command includes ready todos
- the command includes blocked todos
- readiness is not used to exclude items from the result set

### `--all` does not take a value

Given:

- the new flag is available

Then:

- the command surface uses `--all` with no argument
- the presence of `--all` alone changes behavior
- forms such as `--all=true` or `--all false` are not required command forms

### Other filters still work

Given:

- the user combines `--all` with status or priority filters

Then:

- status and priority filters continue to work as before
- `--all` affects only readiness-based filtering

### Usage and help text are updated

Given:

- the feature is implemented

Then:

- usage/help output documents `--all`
- usage/help output no longer documents `--ready <true|false>` for `list`

## Decisions

1. `--all` is a presence-only flag.
2. The default `list` behavior remains ready-only for non-done todos.
3. `--all` disables readiness filtering rather than inverting it.

## Out of scope

- changing the computed `ready` field in JSON output
- changing status filter semantics
- adding a separate flag for blocked-only filtering
