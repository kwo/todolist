---
id: todo-m9yf
title: 'User Story: split cli.go into command-specific files'
status: todo
priority: 4
createdAt: "2026-03-23T19:13:37Z"
lastModified: "2026-03-23T19:16:27Z"
---

# User Story: split `cli.go` into command-specific files

Refactor the CLI package so that command parsing and execution logic are split across multiple files, with command-related code grouped by command.

## User story

As a developer,
I want `pkg/cli/cli.go` to be split into multiple files organized by command,
so that the CLI implementation is easier to navigate, maintain, and extend.

## Goal

This work restructures the `pkg/cli` package to reduce the size and scope of `cli.go` by moving command-specific logic into separate files.

## Desired structure

The exact filenames may vary, but the package should move toward a structure such as:

- `pkg/cli/app.go` — app construction and top-level run flow
- `pkg/cli/parse.go` — root parsing and shared parsing helpers
- `pkg/cli/add.go` — `add` command parsing and execution
- `pkg/cli/list.go` — `list` command parsing and execution
- `pkg/cli/view.go` — `view` command parsing and execution
- `pkg/cli/update.go` — `update` command parsing and execution
- `pkg/cli/delete.go` — `delete` command parsing and execution
- `pkg/cli/init.go` — `init` command parsing and execution
- `pkg/cli/help.go` — help text rendering
- `pkg/cli/json.go` or `pkg/cli/output.go` — shared output helpers

## Recommended organization approach

Use a hybrid structure:

- keep low-level parsing primitives in shared files when they are genuinely reused across commands
- colocate command-specific parsing and execution with the command they belong to

Examples of shared helpers that can remain centralized:

- `parseAssignment`
- `parseExplicitPriority`
- `recognizeStatusValue`
- `recognizePriorityValue`
- `isDigits`
- shared output helpers such as JSON rendering

Examples of helpers that should live with `list` because they are command-specific:

- `parseListCommand`
- `parseStatusFilter`
- `parseExplicitStatusFilter`
- `parsePriorityFilter`
- `recognizeStatusFilter`
- `recognizePriorityFilter`
- list-specific title formatting helpers

Examples of helpers that should live with other commands:

- `parseAddCommand` and `addCommand.Execute` in `add.go`
- `parseUpdateCommand` and `updateCommand.Execute` in `update.go`
- `parseSingleTodoCommand` may remain shared if it is still used by multiple commands such as `view` and `delete`

This approach keeps shared behavior consistent while still making each command easy to find.

## Acceptance criteria

### Command separation

Given:

- the CLI package currently stores most behavior in `pkg/cli/cli.go`

When:

- the refactor is complete

Then:

- command-specific parsing and execution logic are moved into separate files grouped by command
- shared helpers remain in shared files
- behavior remains unchanged

### Test stability

Given:

- the CLI already has test coverage

When:

- the refactor is complete

Then:

- existing tests still pass without requiring behavior changes unrelated to the refactor

### Maintainability

Given:

- a developer needs to change one command

Then:

- the relevant command implementation can be found without scanning a monolithic `cli.go`

## Scope

- refactor of `pkg/cli`
- file and code organization improvements
- no intended user-facing CLI behavior changes

## Open issues

None currently.
