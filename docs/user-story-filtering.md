# User Story: task filtering

Add filtering to task listing.

## User story

As a user,
I want to filter task listings,
so that I can focus on the tasks that matter right now.

## Goal

This work adds filtering to `list` on top of task metadata.

## Command surface

### `list`

```bash
tasklist list [--status <status>] [--priority <priority-filter>]
```

Behavior:

- retain the default compact human-readable list format
- apply filters before rendering output
- `--status <status>` includes only tasks with that status
- `--status !<status>` excludes tasks with that status
- `--priority <n>` includes only tasks with that priority
- `--priority !<n>` excludes tasks with that priority
- `--priority >n` includes only tasks with a numerically greater priority value than `n`
- `--priority <n` includes only tasks with a numerically lower priority value than `n`

Examples:

```bash
tasklist list --status done
tasklist list --status '!done'
tasklist list --priority 1
tasklist list --priority '!3'
tasklist list --priority '>3'
tasklist list --priority '<3'
```

Because priorities are numeric and `1` is the highest priority while `5` is the lowest priority:

- `--priority >3` means priorities `4` and `5`
- `--priority <3` means priorities `1` and `2`

## Acceptance criteria

### Status filtering

Given:

- tasks have a `status` field

When:

- the user runs `tasklist list --status <status>`

Then:

- only tasks with that status are listed

When:

- the user runs `tasklist list --status !<status>`

Then:

- tasks with that status are excluded from the list

### Priority filtering

Given:

- tasks have a `priority` field

When:

- the user runs `tasklist list --priority <n>`

Then:

- only tasks with that priority are listed

When:

- the user runs `tasklist list --priority !<n>`

Then:

- tasks with that priority are excluded from the list

When:

- the user runs `tasklist list --priority >n`

Then:

- only tasks with a numerically greater priority value than `n` are listed

When:

- the user runs `tasklist list --priority <n`

Then:

- only tasks with a numerically lower priority value than `n` are listed

## Scope

- filtering in `list` by status
- filtering in `list` by priority

## Dependencies

- depends on [User Story: task metadata](user-story-metadata.md)

## Open issues

1. Shells may require quoting filters such as `--status '!done'`, `--priority '>3'`, and `--priority '<3'`.
