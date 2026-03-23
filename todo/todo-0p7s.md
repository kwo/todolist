---
id: todo-0p7s
title: 'User Story: todo directory configuration'
status: done
createdAt: "2026-03-23T18:50:56Z"
lastModified: "2026-03-23T18:50:56Z"
---

# User Story: todo directory configuration

Add per-directory todolist configuration via `.todos`.

## User story

As a user,
I want a configuration file in my todo directory,
so that directory-specific behavior such as todo ID prefixes can be customized.

## Goal

This work adds support for reading `.todos` from the todo directory.

## Configuration format

The `.todos` file is a flat `key=value` file.

Supported keys:

- `prefix` — the prefix used when generating todo IDs

If `prefix` is not set, the default prefix is `todo-`.

Example:

```text
prefix=work-
```

## Acceptance criteria

Given:

- the selected todo directory contains a `.todos` file

Then:

- todolist reads `.todos`
- todo ID generation uses the configured `prefix`

Given:

- no `.todos` file exists or `prefix` is not set

Then:

- the default prefix `todo-` is used

Given:

- the `.todos` file contains a malformed line or unsupported key

Then:

- todolist returns an error instead of silently ignoring the invalid configuration

## Scope

- `.todos` file support
- `prefix` setting
- configurable todo ID prefixes

## Dependencies

- naturally paired with [User Story: todo directory selection](user-story-directory-selection.md)
- complemented by [User Story: `init` command](user-story-init.md)

## Configuration parsing notes

- blank lines are ignored
- each non-empty line must use `key=value` format
- unsupported keys are rejected as errors
- an empty `prefix` value falls back to the default prefix `todo-`
