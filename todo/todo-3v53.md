---
id: todo-3v53
title: dependency cycle detection
status: todo
priority: 3
parents:
    - todo-ed12
depends:
    - todo-ed12
createdAt: "2026-03-25T18:46:39Z"
lastModified: "2026-03-27T19:43:36Z"
---

# User Story: dependency cycle detection

Add dependency cycle detection and reporting for todos.

## User story

As a user,
I want to detect dependency cycles,
so that I can find invalid or confusing blocking relationships.

## Goal

This story adds a command that scans the todo graph for cycles and reports them clearly.

## Command surface

```bash
todolist cycles
```

Behavior:

- scans all todos for dependency cycles
- reports whether cycles were found
- reports each detected cycle using todo IDs

## Acceptance criteria

### No cycles

Given:

- a todo set with no dependency cycles

When:

- the user runs the command

Then:

- the CLI reports that no cycles were found
- the command exits successfully

### One or more cycles

Given:

- a todo set with one or more dependency cycles

When:

- the user runs the command

Then:

- the CLI reports each detected cycle using todo IDs
- the command exits non-zero to signal that cycles were found

### Readable reporting

Given:

- a detected cycle

Then:

- the output makes the cycle path understandable to a human reader
- the same cycle is not reported repeatedly in trivial reordered forms

### Multiple cycles

Given:

- a todo set containing multiple distinct cycles

Then:

- the output includes each distinct cycle

### Missing or malformed references

Given:

- a todo references a dependency that is missing, malformed, or otherwise unreadable

Then:

- that condition does not by itself count as a cycle
- the command still analyzes the rest of the graph
- output may report such references separately as warnings

## Decisions

1. The cycle feature is a top-level command.
2. Finding one or more cycles is a detectable problem state, so the command exits non-zero when cycles are present.
3. Cycle reporting should deduplicate equivalent presentations of the same cycle.
4. Missing or malformed references are warnings, not cycles.

## Out of scope

- dependency editing
- readiness computation
- dependency tree visualization
