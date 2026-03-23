---
id: todo-hhsq
title: 'User Story: list output status, priority, and title width'
status: done
priority: 3
createdAt: "2026-03-23T18:58:14Z"
lastModified: "2026-03-23T19:07:18Z"
---

# User Story: list output status, priority, and title width

Add status and priority columns to human-readable `list` output, and limit the displayed title width to 60 characters.

## User story

As a user,
I want `todolist list` to show each todo's status and priority alongside the title,
so that I can scan the current state of my work more quickly.

As a user,
I want long titles in `list` output to be limited to 60 characters,
so that listings stay compact and readable in a terminal.

## Goal

This work updates the default human-readable `list` output to include:

- `id`
- `priority`
- `status`
- `title`, truncated to a maximum display width of 60 characters

## Command surface

### `list`

```bash
todolist list
todolist list done
todolist list 2
todolist list status=done priority=3+
```

The command syntax does not change. Only the default human-readable rendering changes.

## Output format

Human-readable `list` output should render one todo per line and include:

- id
- priority
- status
- title

The first column must be the todo id.

Suggested column order:

```text
<id>\t<priority>\t<status>\t<title>
```

Example:

```text
todo-7k9m\t1\twip\tBuy groceries
todo-2w8x\t5\tdone\tWrite proposal
```

## Title width behavior

For human-readable `list` output:

- the displayed title should be limited to 60 characters maximum
- titles shorter than or equal to 60 characters are shown unchanged
- titles longer than 60 characters are truncated for display only
- if a title is truncated, an ellipsis should be shown at the end of the displayed title
- the total displayed title length, including the ellipsis, must not exceed 60 characters
- the underlying stored title is not modified

Example:

- stored title:
  - `Investigate how to reconcile customer billing exports across regions and vendors`
- displayed title:
  - `Investigate how to reconcile customer billing exports ac...`

## Acceptance criteria

### Status and priority columns

Given:

- the user runs `todolist list`

Then:

- each listed todo includes its id, priority, status, and title in the default human-readable output
- the first column is always the todo id

### Title width limiting

Given:

- a todo title is longer than 60 characters

When:

- the user runs `todolist list`

Then:

- the displayed title is truncated to 60 characters maximum
- if truncation occurs, the displayed title ends with an ellipsis
- the total displayed title length, including the ellipsis, does not exceed 60 characters
- the todo's stored title remains unchanged

### Filtering compatibility

Given:

- the user applies existing list filters

Then:

- filtering behavior remains unchanged
- only the rendering format differs

### JSON compatibility

Given:

- the user runs `todolist list --json`

Then:

- JSON output remains unchanged
- the full title value is preserved in JSON

## Scope

- human-readable `list` output formatting
- id-first column layout in `list`
- status column in `list`
- priority column in `list`
- title display truncation to 60 characters in `list`
- ellipsis for truncated titles within the 60-character maximum

## Open issues

None currently.
