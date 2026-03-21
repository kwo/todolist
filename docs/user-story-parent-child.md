# User Story: parent-child todo grouping

Add parent-child todo grouping to todolist.

## User story

As a user,
I want to group todos under a parent todo,
so that I can organize related work.

## Goal

This work adds parent-child grouping via `parentIds`.

## Todo format additions

This user story extends todo front matter with:

- `parentIds` — optional list of parent todo IDs used for grouping

Example:

```md
---
id: todo-7k9m
title: Buy groceries
status: todo
priority: 2
parentIds:
  - todo-3h7q
  - todo-9p2d
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

## Acceptance criteria

### Parent-child grouping

Given:

- a todo may have zero or more parent todos
- a todo may have many child todos

Then:

- parent-child grouping is represented with `parentIds`
- a todo may belong to multiple parents at the same time
- parent-child grouping is separate from dependency relationships

### Parent deletion behavior

Given:

- a parent todo has child todos

When:

- the user tries to delete the parent todo

Then:

- deletion should warn and fail by default
- the parent cannot be deleted until its child todos are moved out of that parent
- if deletion is invoked with an explicit forced-deletion control such as `force=true`, deletion may proceed and should also delete all child todos

## Scope

- `parentIds` support for parent-child grouping
- guarded deletion behavior for parent todos with children

## Open issues

1. The CLI behavior for creating, updating, or listing parent-child relationships beyond `parentIds` storage is not yet specified.
