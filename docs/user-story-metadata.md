# User Story: todo metadata

Add todo metadata fields beyond the MVP core fields.

## User story

As a user,
I want todos to carry status and priority metadata,
so that I can better track progress and urgency.

## Goal

This work adds:

- `status`
- `priority`
- support for setting and updating those fields through the CLI

## Todo format additions

This user story extends todo front matter with:

- `status`
- `priority`

Example:

```md
---
id: todo-7k9m
title: Buy groceries
status: todo
priority: 2
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

## Command surface

### `add`

```bash
todolist add <title> [<status>] [<priority>]
todolist add title=<title> [status=<status>] [priority=<priority>]
```

### `update`

```bash
todolist update <todo> [<title>] [<status>] [<priority>]
todolist update <todo> [title=<title>] [status=<status>] [priority=<priority>]
```

### `view`

```bash
todolist view <todo>
```

Behavior:

- `view` should show all todo metadata, including `status` and `priority`

## Acceptance criteria

### Status

Given:

- a todo may have a status

Then:

- supported values are `todo`, `wip`, and `done`, matching the shared status values documented in `README.md`
- if status is omitted, it defaults to `todo`
- `add` may set status using inferred values or explicit `status=<status>` notation
- `update` may replace status using inferred values or explicit `status=<status>` notation

### Priority

Given:

- a todo may have a priority

Then:

- supported values are `1`, `2`, `3`, `4`, and `5`, matching the shared priority values documented in `README.md`
- if the `priority` attribute is omitted, it defaults to `5`
- `add` may set priority using inferred values or explicit `priority=<priority>` notation
- `update` may replace priority using inferred values or explicit `priority=<priority>` notation

### Timestamps

Given:

- a todo is updated through `update`

Then:

- `lastModified` should be updated automatically

## Scope

- `status` metadata
- `priority` metadata
- `add` support for inferred `status` and `priority` values, plus explicit `status=...` and `priority=...` notation
- `update` support for inferred `title`, `status`, and `priority` values, plus explicit `title=...`, `status=...`, and `priority=...` notation
- `view` output includes metadata values

## Open issues

1. Validation and error messages for invalid `status` and `priority` values should be made explicit.
