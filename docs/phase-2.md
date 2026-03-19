# Phase 2: Task relationships and dependencies

Phase 2 adds parent-child grouping, dependency tracking, readiness, and dependency inspection.

## Commands

- `dep` — manage task dependencies

Dependency commands inherit the root CLI's global options, including `-d, --directory`. Because `--directory` is global, it should appear before the subcommand. Example:

```bash
tasklist -d ./work-tasks dep list task-7k9m
```

## Task Format Additions

Phase 2 extends the front matter with:

- `parentId` — optional parent task id used for grouping
- `dependsOn` — optional list of task ids this task depends on

Phase 2 also adds a computed field:

- `ready` — true when all dependencies are completed; false when one or more dependencies are not yet completed

## Task Relationships

### Parent-child grouping

A task may have:

- one optional parent task
- many child tasks

This parent-child relationship is for grouping tasks.

If a parent task has child tasks, deleting the parent should emit a warning and fail by default. The parent cannot be deleted until its child tasks are moved outside of that parent, unless deletion is invoked with `-f` or `--force`. With the force option, deleting the parent may proceed and should also delete all of its child tasks.

### Dependencies

Dependencies are separate from parent-child grouping.

Dependencies should be stored in `dependsOn` as a list of task ids.

A task also has a computed `ready` field. `ready` is true only when all dependencies are completed.

The CLI will not prevent circular dependencies automatically, but it should provide commands to inspect and report them.

## Dependency Commands

Phase 2 should provide a `dep` command with these subcommands:

- `dep add <task> <depends-on>`
  - Add a dependency
  - Meaning: `<task>` depends on `<depends-on>`

- `dep remove <task> <depends-on>`
  - Remove a dependency
  - Alias: `dep rm`

- `dep list <task>`
  - List dependencies for a task
  - Supports directions:
    - `--direction down` = what this task depends on
    - `--direction up` = what depends on this task
    - `--direction both`

- `dep tree <task>`
  - Show a dependency tree rooted at the task
  - Options include:
    - `--direction`
    - `--max-depth`
    - `--format text|mermaid`

- `dep cycles`
  - Detect and report dependency cycles

## Scope

- `parentId` support for parent-child grouping
- `dependsOn` support in task metadata
- computed `ready` field
- dependency traversal and visualization
- dependency cycle detection

## Open Issues

1. Delete behavior with dependencies is unspecified. If you delete a task that other tasks depend on (via `dependsOn`), what happens to those references? Only parent-child deletion is addressed.
2. `ready` definition of "completed" is ambiguous. It says "all dependencies are completed" but doesn't say what "completed" means. Presumably `status: done`, but this should be explicit.
