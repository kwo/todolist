# User Story: task directory configuration

Add per-directory tasklist configuration via `.tasks`.

## User story

As a user,
I want a configuration file in my task directory,
so that directory-specific behavior such as task ID prefixes can be customized.

## Goal

This work adds support for reading `.tasks` from the task directory.

## Configuration format

The `.tasks` file is a flat `key=value` file.

Supported keys:

- `prefix` — the prefix used when generating task IDs

If `prefix` is not set, the default prefix is `task-`.

Example:

```text
prefix=work-
```

## Acceptance criteria

Given:

- the selected task directory contains a `.tasks` file

Then:

- tasklist reads `.tasks`
- task ID generation uses the configured `prefix`

Given:

- no `.tasks` file exists or `prefix` is not set

Then:

- the default prefix `task-` is used

## Scope

- `.tasks` file support
- `prefix` setting
- configurable task ID prefixes

## Dependencies

- naturally paired with [User Story: task directory selection](user-story-directory-selection.md)
- complemented by [User Story: `init` command](user-story-init.md)

## Open issues

1. The behavior for malformed `.tasks` files should be specified explicitly.
