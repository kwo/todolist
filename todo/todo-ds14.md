---
id: todo-ds14
title: 'User Story: parent-child todo grouping'
status: done
priority: 1
createdAt: "2026-03-23T18:50:57Z"
lastModified: "2026-03-25T18:30:41Z"
---

# User Story: parent-child todo grouping

Add parent-child todo grouping to `todolist`.

## User story

As a user,
I want to group todos under one or more parent todos,
so that I can organize related work without using dependency relationships.

## Goal

This work adds first-class parent-child grouping via `parents` in todo front matter and defines the minimum CLI behavior needed to create, inspect, update, and safely delete grouped todos.

## Problem statement

Today todos are flat. Users can model dependencies, but they cannot represent hierarchical organization such as epics, projects, or grouped subtasks. This story introduces a lightweight grouping model where a todo may belong to zero or more parents.

## Todo format additions

This story extends todo front matter with:

- `parents` — optional list of parent todo IDs used only for grouping

Example:

```md
---
id: todo-7k9m
title: Buy groceries
status: todo
priority: 2
parents:
  - todo-3h7q
  - todo-9p2d
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

## Domain rules

- A todo may have zero or more parent todos.
- A todo may have zero or more child todos.
- A todo may belong to multiple parents at the same time.
- `parents` is for grouping only and must not change dependency behavior.
- Parent relationships must reference existing todo IDs.
- A todo must not list itself as a parent.
- Stored `parents` should be unique and stable in output order.
- Omitting `parents` and providing an empty `parents` list are treated equivalently.

## CLI behavior

### Create

Users can create a todo with zero or more parent IDs.

Acceptance criteria:

- `todolist add` accepts repeated `--parent <todo-id>` flags.
- Parent IDs are validated before the todo is written.
- Creating a todo fails with a clear error if any parent ID does not exist.
- Creating a todo fails with a clear error if duplicate parent IDs are provided.

Example:

```bash
todolist add --title "Buy groceries" --parent todo-3h7q --parent todo-9p2d
```

### Update

Users can add and remove parent relationships on an existing todo.

Acceptance criteria:

- `todolist update` accepts repeated `--parent <todo-id>` flags to add parents.
- `todolist update` accepts repeated `--parent <todo-id>!` flags to remove parents.
- Parent update operations are applied in flag order.
- Updating a todo validates all referenced parent IDs.
- Updating fails if any added parent does not exist.
- Updating fails if the todo would reference itself as a parent.
- Updating fails with a clear error if the same update request contains conflicting operations for the same parent.
- Updating fails with a clear error if a remove operation targets a parent that is not currently assigned.
- Updating parent relationships does not modify dependency relationships.
- Updating can remove all parent relationships by removing each current parent with the `!` suffix form.

Examples:

```bash
todolist update todo-7k9m --parent todo-3h7q --parent todo-9p2d
todolist update todo-7k9m --parent todo-3h7q!
todolist update todo-7k9m --parent todo-3h7q --parent todo-9p2d!
```

### Read and list

Users can inspect parent relationships.

Acceptance criteria:

- `todolist view --json` includes `parents` when present.
- `todolist list --json` includes `parents` for each todo.
- Text `todolist view` output displays parents in a human-friendly section, not only as raw front matter.
- The human-friendly section should appear even though `parents` is also stored in front matter.
- In the human-friendly parents section, each parent is displayed across multiple lines rather than on a single line.
- Each parent entry includes both the todo ID and the title.
- Non-JSON list output includes a parent IDs column as the last column.
- The parent IDs column displays only the first parent ID.
- If a todo has multiple parent IDs, the column shows the first parent ID followed by an ellipsis.
- Non-JSON list output remains readable and does not need to display the full hierarchy beyond that column.

## Deletion behavior

When a parent todo is deleted:

- deletion succeeds even if other todos reference it in `parents`
- no child todos are deleted
- the deleted todo's ID is removed from the `parents` list of every affected child todo
- child todos remain otherwise unchanged
- removing the deleted parent from child todos does not modify dependency relationships
- cleanup is applied consistently across all todos before the delete operation completes

Example:

- `todo-parent` is a parent of `todo-a` and `todo-b`
- deleting `todo-parent` removes `todo-parent` from `todo-a.parents` and `todo-b.parents`
- `todo-a` and `todo-b` continue to exist

## Data integrity expectations

- Reading existing todos without `parents` continues to work.
- Existing todo files remain valid without migration.
- Serialization preserves backward compatibility for todos that do not use grouping.
- Invalid on-disk `parents` should surface a clear validation error when parsed.
- Deleting a todo must not leave dangling parent references in remaining todos.

## Out of scope

- Visual tree rendering in text output
- New dependency semantics
- Automatic parent status rollups from children
- Sorting or filtering by parent unless needed to support the minimum parent feature

## Implementation notes

- Keep the data model change minimal: add `parents` to the todo schema and validation logic.
- Reuse existing todo lookup paths when validating parents.
- Use repeated `--parent` flags rather than a comma-separated value.
- In `update`, treat `--parent <id>` as add and `--parent <id>!` as remove.
- Delete logic should find all todos that reference the deleted ID in `parents` and remove that reference before completing the deletion.
- Update CLI usage docs once the flag shape is finalized.

