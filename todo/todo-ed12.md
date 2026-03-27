---
id: todo-ed12
title: todo dependencies
status: done
priority: 2
createdAt: "2026-03-23T18:50:57Z"
lastModified: "2026-03-27T19:54:11Z"
---

# User Story: todo dependencies

Add dependency tracking and readiness reporting to todolist.

## User stories

### Dependencies

As a user,
I want to record that one todo depends on another,
so that I can understand what is blocked and what is ready to work on.

### Readiness

As a user,
I want each todo to expose whether it is ready,
so that I can quickly identify work that is unblocked.

## Goal

This work adds:

- dependency tracking via `depends`
- computed readiness via `ready`
- dependency editing through the existing `add` and `update` commands
- dependency visibility in `list` output

## Command surface

This story does **not** add a new `dep` command.

Instead, dependencies are managed through flags on existing commands.

### `add`

```bash
todolist add -t <title> [--depends <todo-id> ...]
```

Behavior:

- `--depends <todo-id>` adds a dependency when creating a todo
- meaning: the new todo depends on `<todo-id>`
- the flag may be repeated

### `update`

```bash
todolist update <todo-id> [--depends <todo-id>|<todo-id>! ...]
```

Behavior:

- `--depends <todo-id>` adds a dependency
- `--depends <todo-id>!` removes a dependency
- the flag may be repeated
- `update` must still change at least one field, including dependency edits

### `list`

`list` output adds a dependency column.

Behavior:

- the dependency column behaves like the existing parent column
- when a todo has one dependency, the column shows that dependency ID
- when a todo has multiple dependencies, the column shows the first dependency ID followed by `,...`
- when a todo has no dependencies, the column is empty

Tree visualization and cycle inspection are intentionally out of scope for this story and will be tracked in separate todos.

## Todo format additions

This user story extends todo front matter with:

- `depends` — optional list of todo IDs this todo depends on

This user story also adds a computed field:

- `ready` — true when all dependencies have `status: done`; false when one or more dependencies are not done

Example:

```md
---
id: todo-7k9m
title: Buy groceries
status: todo
priority: 2
depends:
  - todo-2w8x
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

## Acceptance criteria

### Dependency storage

Given:

- one todo depends on another

Then:

- the relationship is stored in `depends`
- the meaning is: `<todo>` depends on `<depends-on>`

### Add command support

Given:

- a user creates a todo with one or more `--depends <todo-id>` flags

Then:

- each referenced todo ID is stored in `depends`

### Update command support

Given:

- a user updates a todo with `--depends <todo-id>`

Then:

- that dependency is added to `depends`

Given:

- a user updates a todo with `--depends <todo-id>!`

Then:

- that dependency is removed from `depends`

### Edge cases

Given:

- a user supplies the same dependency more than once during `add` or `update`

Then:

- the dependency is stored at most once
- duplicate additions do not create duplicate entries in `depends`

Given:

- a user tries to remove a dependency that is not currently present

Then:

- the command fails

Given:

- a user tries to make a todo depend on itself

Then:

- the command fails
- the todo is not modified

Given:

- a user references a dependency ID that does not exist

Then:

- the command fails
- the todo is not modified

### Readiness

Given:

- a todo has zero or more dependencies

Then:

- `ready` is computed, not stored
- `ready` is true only when all dependencies have `status: done`
- `ready` is false when one or more dependencies have any other status, such as `todo` or `wip`
- if a dependency is missing a status, it is treated as `todo`
- if a dependency is malformed, invalid, or otherwise unreadable, it is treated as not done, so `ready` is false

### Delete behavior

Given:

- a todo is deleted
- one or more remaining todos reference it in `depends`

Then:

- the delete still succeeds
- the deleted todo ID is removed from every remaining todo's `depends` list
- each affected todo's `lastModified` is updated
- readiness is recomputed from the remaining dependencies

### List output

Given:

- a user runs `todolist list`

Then:

- text output includes a dependency column
- the dependency column behaves like the parent column
- JSON output includes `depends`
- JSON output includes computed `ready`

## Scope

- `depends` support in todo metadata
- computed `ready` field
- `--depends` support in `add`
- `--depends` add/remove support in `update`
- dependency visibility in `list`

## Out of scope

The following capabilities will be split into separate todos:

- dependency tree traversal and visualization
- dependency cycle detection and reporting

## Decisions

1. There is no `dep` command in this story; dependency editing happens through `add` and `update`.
2. `list` gains a dependency column that mirrors the parent column behavior.
3. Tree and cycle functionality are not part of this story and will be tracked separately.
4. If a todo is deleted and other todos depend on it via `depends`, those references are removed automatically as part of deletion.
5. A missing dependency status is treated as the default status `todo`.
6. A dependency whose content is malformed, invalid, or otherwise unreadable is treated as not done when computing `ready`.
7. Duplicate dependencies are deduplicated.
8. Removing a dependency that is not present fails.
9. Self-dependencies are rejected.
10. Referencing a nonexistent dependency ID is rejected.
