---
id: todo-me9x
title: dependency tree visualization
status: todo
priority: 3
parents:
    - todo-ed12
createdAt: "2026-03-25T18:46:29Z"
lastModified: "2026-03-25T18:50:17Z"
---

# User Story: dependency tree visualization

Add dependency tree inspection for todos.

## User story

As a user,
I want to view a dependency tree for a todo,
so that I can understand upstream blockers and downstream dependents.

## Goal

This story adds a command for traversing dependency relationships from a root todo and rendering the result as either text or Mermaid.

## Command surface

```bash
todolist tree <todo> [direction=down|up|both] [max-depth=<n>] [format=text|mermaid]
```

Behavior:

- `<todo>` is the root of the traversal
- `direction=down` traverses `dependsOn` edges from the root to its dependencies
- `direction=up` traverses reverse dependency edges from the root to dependents
- `direction=both` includes both upstream and downstream relationships
- `max-depth=<n>` limits traversal depth from the root
- `format=text` prints a human-readable tree
- `format=mermaid` prints Mermaid output suitable for visualization

## Acceptance criteria

### Root lookup

Given:

- the user supplies a root todo ID

Then:

- the command loads that todo as the tree root
- if the root todo does not exist, the command fails

### Traversal directions

Given:

- a root todo with dependencies

When:

- the user requests `direction=down`

Then:

- the output shows the todos that the root depends on

Given:

- a root todo with dependents

When:

- the user requests `direction=up`

Then:

- the output shows the todos that depend on the root

Given:

- a root todo with both dependencies and dependents

When:

- the user requests `direction=both`

Then:

- the output includes both upstream and downstream relationships

### Depth limiting

Given:

- a root todo and a traversal depth limit

When:

- the user supplies `max-depth=<n>`

Then:

- traversal stops after `n` levels from the root
- the root is depth 0

### Text output

Given:

- the user requests `format=text`

Then:

- the CLI prints a readable tree
- each entry identifies the todo by ID
- each entry includes enough information to distinguish siblings, such as title
- the root todo is clearly indicated

### Mermaid output

Given:

- the user requests `format=mermaid`

Then:

- the CLI prints valid Mermaid graph content
- edges reflect the selected traversal direction
- the root todo appears in the graph

### Safe traversal in the presence of cycles

Given:

- the traversed graph contains a cycle

Then:

- the command does not recurse forever
- repeated visits are handled safely
- output indicates recursion was cut off when necessary

### Missing referenced todos

Given:

- a todo references a dependency ID that cannot be loaded

Then:

- the command still returns a useful result for the reachable graph
- the missing reference is represented clearly in output

## Decisions

1. The tree feature is a top-level command, not part of `add`, `update`, or `list`.
2. The root todo is depth 0 for `max-depth` calculations.
3. Cycles do not cause failure by default; traversal remains safe and finite.
4. Missing referenced todos should be shown clearly rather than causing the entire command to fail, unless the root todo itself is missing.

## Out of scope

- dependency editing
- readiness computation
- full cycle detection and reporting policy beyond safe traversal behavior
