# User Story: parent-child task grouping

Add parent-child task grouping to tasklist.

## User story

As a user,
I want to group tasks under a parent task,
so that I can organize related work.

## Goal

This work adds parent-child grouping via `parentId`.

## Task format additions

This user story extends task front matter with:

- `parentId` — optional parent task ID used for grouping

Example:

```md
---
id: task-7k9m
title: Buy groceries
status: todo
priority: 2
parentId: task-3h7q
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

## Acceptance criteria

### Parent-child grouping

Given:

- a task may have one optional parent task
- a task may have many child tasks

Then:

- parent-child grouping is represented with `parentId`
- parent-child grouping is separate from dependency relationships

### Parent deletion behavior

Given:

- a parent task has child tasks

When:

- the user tries to delete the parent task

Then:

- deletion should warn and fail by default
- the parent cannot be deleted until its child tasks are moved out of that parent
- if deletion is invoked with `-f` or `--force`, deletion may proceed and should also delete all child tasks

## Scope

- `parentId` support for parent-child grouping
- guarded deletion behavior for parent tasks with children

## Open issues

1. The CLI behavior for creating, updating, or listing parent-child relationships beyond `parentId` storage is not yet specified.
