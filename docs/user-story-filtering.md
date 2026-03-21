# User Story: todo filtering

Add filtering to todo listing.

## User story

As a user,
I want to filter todo listings,
so that I can focus on the todos that matter right now.

## Goal

This work adds filtering to `list` on top of todo metadata.

## Command surface

### `list`

```bash
todolist list [<status-filter>] [<priority-filter>]
todolist list [status=<status-filter>] [priority=<priority-filter>]
```

Behavior:

- retain the default compact human-readable list format
- apply filters before rendering output
- a positional status filter includes only todos with that status
- a positional `!<status>` filter excludes todos with that status
- a positional priority filter `<n>` includes only todos with that priority
- a positional priority filter `.<n>` excludes todos with that priority
- a positional priority filter `+<n>` includes only todos with a numerically greater priority value than `n`
- a positional priority filter `-<n>` includes only todos with a numerically lower priority value than `n`
- explicit `status=<status-filter>` and `priority=<priority-filter>` notation may also be used

Examples:

```bash
todolist list done
todolist list '!done'
todolist list 1
todolist list .3
todolist list +3
todolist list priority=-3
todolist list status=done
todolist list priority=+3
```

Because priorities are numeric and `1` is the highest priority while `5` is the lowest priority:

- `+0` means priorities `1` through `5`
- `+3` means priorities `4` and `5`
- `-3` means priorities `1` and `2`

## Acceptance criteria

### Status filtering

Given:

- todos have a `status` field

When:

- the user runs `todolist list <status>` or `todolist list status=<status>`

Then:

- only todos with that status are listed

When:

- the user runs `todolist list !<status>` or `todolist list status=!<status>`

Then:

- todos with that status are excluded from the list

### Priority filtering

Given:

- todos have a `priority` field

When:

- the user runs `todolist list <n>` or `todolist list priority=<n>`

Then:

- only todos with that priority are listed

When:

- the user runs `todolist list .<n>` or `todolist list priority=.<n>`

Then:

- todos with that priority are excluded from the list

When:

- the user runs `todolist list +<n>` or `todolist list priority=+<n>`

Then:

- only todos with a numerically greater priority value than `n` are listed

When:

- the user runs `todolist list -<n>` or `todolist list priority=-<n>`

Then:

- only todos with a numerically lower priority value than `n` are listed

## Scope

- filtering in `list` by status
- filtering in `list` by priority

## Dependencies

- depends on [User Story: todo metadata](user-story-metadata.md)

## Open issues

1. Shells may require quoting filters such as `!done`.
2. It should be specified whether a bare value like `-3` is always accepted as a priority filter or whether `priority=-3` should be the preferred explicit form.
