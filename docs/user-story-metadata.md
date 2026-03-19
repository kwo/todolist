# User Story: task metadata

Add task metadata fields beyond the MVP core fields.

## User story

As a user,
I want tasks to carry status and priority metadata,
so that I can better track progress and urgency.

## Goal

This work adds:

- `status`
- `priority`
- support for setting and updating those fields through the CLI

## Task format additions

This user story extends task front matter with:

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

## Command surface

### `add`

```bash
tasklist add <title> [--status <status>] [--priority <priority>]
```

### `update`

```bash
tasklist update <task> [--title <title>] [--status <status>] [--priority <priority>]
```

### `view`

```bash
tasklist view <task>
```

Behavior:

- `view` should show all task metadata, including `status` and `priority`

## Acceptance criteria

### Status

Given:

- a task may have a status

Then:

- supported values are the shared status values documented in `README.md`
- if status is omitted, it defaults to `todo`
- `add` may set status
- `update` may replace status

### Priority

Given:

- a task may have a priority

Then:

- supported values are the shared priority values documented in `README.md`
- if priority is omitted, it defaults to `0`
- `add` may set priority
- `update` may replace priority

### Timestamps

Given:

- a task is updated through `update`

Then:

- `lastModified` should be updated automatically

## Scope

- `status` metadata
- `priority` metadata
- `add` support for `--status` and `--priority`
- `update` support for `--status` and `--priority`
- `view` output includes metadata values

## Open issues

1. Validation and error messages for invalid `status` and `priority` values should be made explicit.
