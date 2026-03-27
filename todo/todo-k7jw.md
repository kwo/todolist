---
id: todo-k7jw
title: default list filter shows only ready items
status: done
priority: 5
parents:
    - todo-ed12
depends:
    - todo-ed12
createdAt: "2026-03-27T19:53:45Z"
lastModified: "2026-03-27T20:27:19Z"
---

# User Story: ready-only default list filter

Update `todolist list` so readiness becomes a first-class filter instead of a text-only display column.

## User stories

### Default ready filtering

As a user,
I want `todolist list` to show only ready work by default,
so that the command focuses on items I can act on immediately.

### Ready filter control

As a user,
I want to explicitly filter on readiness,
so that I can switch between ready and blocked work without changing other filters.

### Documentation sync

As a user,
I want the usage guide to match the CLI behavior,
so that the documented defaults and flags stay accurate.

## Goal

This work changes `list` to:

- default to showing only todos where computed `ready` is `true`
- add a new `--ready` flag to filter on computed readiness
- remove the `ready` column from text `list` output
- keep computed `ready` in JSON output
- update `USAGE.md` to describe the new default behavior, the `--ready` flag, and the revised text output columns

## Command surface

### `list`

```bash
todolist list [--status <status-filter>] [--priority <priority-filter>] [--ready <true|false>]
```

Behavior:

- `--ready true` shows only todos whose computed `ready` is `true`
- `--ready false` shows only todos whose computed `ready` is `false`
- when `--ready` is omitted, the default is `true`
- readiness filtering is applied in addition to existing status and priority filters

## Output behavior

### Text output

Text `list` output removes the `ready` column.

Behavior:

- text output returns the same columns as before readiness was displayed
- the dependency column remains present
- the column order is:
  - `<id>`
  - `<priority>`
  - `<status>`
  - `<title>`
  - `<first-parent-id>`
  - `<first-dependency-id>`

### JSON output

JSON `list` output continues to include computed `ready`.

Behavior:

- `ready` remains computed, not stored
- JSON output includes `depends`
- JSON output includes computed `ready`

## Acceptance criteria

### Default filtering

Given:

- one ready todo
- one blocked todo

When:

- the user runs `todolist list`

Then:

- the ready todo is included
- the blocked todo is excluded

### Explicit ready filtering

Given:

- one ready todo
- one blocked todo

When:

- the user runs `todolist list --ready true`

Then:

- only the ready todo items are included

Given:

- one ready todo
- one blocked todo

When:

- the user runs `todolist list --ready false`

Then:

- only the blocked todo items are included

### Filter composition

Given:

- todos with different combinations of status, priority, and readiness

When:

- the user runs `todolist list` with `--ready` and other filters

Then:

- the result includes only todos that satisfy all supplied filters

### Text list output

Given:

- a user runs `todolist list`

Then:

- text output does not include a `ready` column
- text output still includes the dependency column

### JSON list output

Given:

- a user runs `todolist list --json`

Then:

- JSON output includes computed `ready`
- JSON output includes `depends`

### Usage guide

Given:

- the `list` command behavior changes in this story

Then:

- `USAGE.md` is updated to describe the new default readiness filter
- `USAGE.md` documents the `--ready <true|false>` flag
- `USAGE.md` reflects the text `list` output columns after the `ready` column is removed
- `USAGE.md` continues to describe JSON output as including computed `ready`

## Scope

- default `list` readiness filter of `true`
- explicit `--ready <true|false>` filtering in `list`
- removal of the text-only `ready` column from `list`
- preservation of computed `ready` in JSON output
- `USAGE.md` updates for the revised `list` behavior

## Out of scope

- changing how `ready` is computed
- changing `view --json` readiness output
- adding readiness filtering to commands other than `list`

## Decisions

1. `list` defaults to `--ready true` when the flag is not supplied.
2. `--ready false` shows blocked items, meaning todos whose computed `ready` is `false`.
3. The `ready` column is removed from text `list` output only.
4. JSON `list` output keeps computed `ready`.
5. Readiness filtering composes with existing status and priority filters.
6. `USAGE.md` must be updated as part of this story.
