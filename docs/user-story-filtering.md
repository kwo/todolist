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
tasklist list [<status-filter>] [<priority-filter>]
tasklist list [status=<status-filter>] [priority=<priority-filter>]
```

Behavior:

- retain the default compact human-readable list format
- apply filters before rendering output
- a positional status filter includes only tasks with that status
- a positional `!<status>` filter excludes tasks with that status
- a positional priority filter `<n>` includes only tasks with that priority
- a positional priority filter `.<n>` excludes tasks with that priority
- a positional priority filter `+<n>` includes only tasks with a numerically greater priority value than `n`
- a positional priority filter `-<n>` includes only tasks with a numerically lower priority value than `n`
- explicit `status=<status-filter>` and `priority=<priority-filter>` notation may also be used

Examples:

```bash
tasklist list done
tasklist list '!done'
tasklist list 1
tasklist list .3
tasklist list +3
tasklist list priority=-3
tasklist list status=done
tasklist list priority=+3
```

Because priorities are numeric and `1` is the highest priority while `5` is the lowest priority:

- `+0` means priorities `1` through `5`
- `+3` means priorities `4` and `5`
- `-3` means priorities `1` and `2`

## Acceptance criteria

### Status filtering

Given:

- tasks have a `status` field

When:

- the user runs `tasklist list <status>` or `tasklist list status=<status>`

Then:

- only tasks with that status are listed

When:

- the user runs `tasklist list !<status>` or `tasklist list status=!<status>`

Then:

- tasks with that status are excluded from the list

### Priority filtering

Given:

- tasks have a `priority` field

When:

- the user runs `tasklist list <n>` or `tasklist list priority=<n>`

Then:

- only tasks with that priority are listed

When:

- the user runs `tasklist list .<n>` or `tasklist list priority=.<n>`

Then:

- tasks with that priority are excluded from the list

When:

- the user runs `tasklist list +<n>` or `tasklist list priority=+<n>`

Then:

- only tasks with a numerically greater priority value than `n` are listed

When:

- the user runs `tasklist list -<n>` or `tasklist list priority=-<n>`

Then:

- only tasks with a numerically lower priority value than `n` are listed

## Scope

- filtering in `list` by status
- filtering in `list` by priority

## Dependencies

- depends on [User Story: task metadata](user-story-metadata.md)

## Open issues

1. Shells may require quoting filters such as `!done`.
2. It should be specified whether a bare value like `-3` is always accepted as a priority filter or whether `priority=-3` should be the preferred explicit form.
