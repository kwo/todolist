---
id: todo-gjza
title: sort todos in list by priority, title
status: todo
priority: 1
createdAt: "2026-03-25T18:51:56Z"
lastModified: "2026-03-25T18:51:56Z"
---

# Sort todos in list output by priority, then title

Update the `list` command so todos are ordered deterministically by priority first and title second.

## User story

As a user,
I want `todolist list` to sort todos by priority and then title,
so that the most important work appears first in a predictable order.

## Acceptance criteria

### Primary sort

Given:

- multiple todos with different priorities

When:

- the user runs `todolist list`

Then:

- todos are sorted by priority ascending
- priority `1` appears before priority `2`, and so on

### Secondary sort

Given:

- multiple todos with the same priority

When:

- the user runs `todolist list`

Then:

- those todos are sorted by title ascending

### JSON output

Given:

- the user runs `todolist list --json`

Then:

- the returned todos follow the same ordering rules

## Decisions

1. Sorting is by priority ascending, then title ascending.
2. The same ordering applies to text and JSON list output.
