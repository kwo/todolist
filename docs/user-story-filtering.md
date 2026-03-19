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
tasklist list [--priority <priority>]
```

Behavior:

- retain the default compact human-readable list format
- apply filters before rendering output

## Acceptance criteria

### Priority filtering

Given:

- tasks have a `priority` field

When:

- the user runs `tasklist list --priority <priority>`

Then:

- only tasks with that priority are listed

## Scope

- filtering in `list` by priority

## Dependencies

- depends on [User Story: task metadata](user-story-metadata.md)

## Open issues

1. Additional filters such as `--status` may be desirable later but are not yet specified.
