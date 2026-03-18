# Task List

Task List is a local-first CLI for managing tasks stored as Markdown files with YAML front matter.

Each task is a plain text file in a directory, making tasks easy to inspect, edit, back up, and version with Git.

## Goals

- store tasks as Markdown files
- use front matter for task metadata
- keep the format human-readable and manually editable
- provide a simple CLI for common task operations
- support machine-readable output with a global `--json` option

## Commands

- `list` — list tasks
- `add` — create a task
- `view` — show a task
- `update` — update task metadata or body via CLI flags
- `delete` — permanently remove a task

A global `--json` option should output command results as JSON.

## Task Format

Each task is a Markdown file with YAML front matter and an optional Markdown body.

Example:

```md
---
id: task-001
title: Buy groceries
status: todo
priority: 2
parentId: errands
dependsOn:
  - task-000
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

Front matter fields:

- `id` — unique task identifier
- `title`
- `status`
- `priority`
- `parentId` — optional parent task id
- `dependsOn` — optional list of task ids this task depends on
- `createdAt`
- `lastModified`

## Status Values

- `todo`
- `in_progress`
- `done`

If status is omitted, it should default to `todo`.

## Priority Values

- `0` — lowest priority
- `1`
- `2`
- `3` — most urgent

## Task Relationships

A task may have:

- one optional parent task
- many child tasks

This parent-child relationship is for grouping tasks. It is separate from dependencies.

Dependencies should be stored in `dependsOn` as a list of task ids.

## File Naming

Task files should be named:

- `<id>.md`

## Principles

- local-first
- plain text
- scriptable
- git-friendly
- easy to edit manually
