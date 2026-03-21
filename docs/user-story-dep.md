# User Story: task dependencies

Add dependency management and dependency inspection to tasklist.

## User stories

### Dependencies

As a user,
I want to record that one task depends on another,
so that I can understand what is blocked and what is ready to work on.

### Dependency inspection

As a user,
I want to inspect dependency relationships,
so that I can see what a task depends on, what depends on it, and whether there are cycles.

## Goal

This work adds:

- dependency tracking via `dependsOn`
- computed readiness via `ready`
- CLI commands for inspecting and managing dependency relationships

## Command surface

### Dependency command

```bash
tasklist dep <subcommand>
```

Dependency commands inherit the root CLI's global options, including `-d, --directory`. In the command-first CLI, those global options appear after the root command. Example:

```bash
tasklist dep -d ./work-tasks list task-7k9m
```

## Task format additions

This user story extends task front matter with:

- `dependsOn` — optional list of task IDs this task depends on

This user story also adds a computed field:

- `ready` — true when all dependencies have `status: done`; false when one or more dependencies have any other status, such as `todo` or `wip`

Example:

```md
---
id: task-7k9m
title: Buy groceries
status: todo
priority: 2
dependsOn:
  - task-2w8x
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

## Acceptance criteria

### Dependency storage

Given:

- one task depends on another

Then:

- the relationship is stored in `dependsOn`
- the meaning is: `<task>` depends on `<depends-on>`

### Readiness

Given:

- a task has zero or more dependencies

Then:

- `ready` is computed, not stored
- `ready` is true only when all dependencies have `status: done`
- `ready` is false when one or more dependencies have any other status, such as `todo` or `wip`

### Cycle handling

Given:

- dependencies may form a cycle

Then:

- the CLI does not need to prevent circular dependencies automatically
- the CLI must provide a way to inspect and report dependency cycles

## Dependency commands

### `dep add`

```bash
tasklist dep add <task> <depends-on>
```

Behavior:

- add a dependency
- meaning: `<task>` depends on `<depends-on>`

### `dep remove`

```bash
tasklist dep remove <task> <depends-on>
```

Alias:

- `dep rm`

Behavior:

- remove a dependency

### `dep list`

```bash
tasklist dep list <task> [direction=down|up|both]
```

Behavior:

- list dependencies for a task
- `direction=down` = what this task depends on
- `direction=up` = what depends on this task
- `direction=both` = both directions

### `dep tree`

```bash
tasklist dep tree <task> [direction=down|up|both] [max-depth=<n>] [format=text|mermaid]
```

Behavior:

- show a dependency tree rooted at the task
- support direction control
- support maximum traversal depth
- support text and Mermaid output

### `dep cycles`

```bash
tasklist dep cycles
```

Behavior:

- detect and report dependency cycles

## Scope

- `dependsOn` support in task metadata
- computed `ready` field
- dependency traversal and visualization
- dependency cycle detection

## Open issues

1. Delete behavior with dependencies is unspecified. If a task is deleted and other tasks depend on it via `dependsOn`, what should happen to those references?
2. The handling of dependencies whose status field is missing should be made explicit. Presumably they should be treated as `status: todo`.
