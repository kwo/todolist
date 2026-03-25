---
id: todo-xdhm
title: split CLI unit tests into command-specific files
status: done
priority: 3
createdAt: "2026-03-25T19:00:00Z"
lastModified: "2026-03-25T19:19:07Z"
---

# Split CLI unit tests into command-specific test files

Refactor the CLI test suite so test files are organized to match the command implementation files.

## User story

As a maintainer,
I want CLI unit tests split into separate test files by command,
so that the test suite is easier to navigate and maintain as the CLI grows.

## Goal

This work reorganizes `pkg/cli` tests so command-related tests live alongside the corresponding command area rather than accumulating in a single large test file.

## Acceptance criteria

### File organization

Given:

- CLI commands are implemented in separate files

Then:

- tests for each command are moved into separate test files that mirror that structure
- shared helpers may remain in a common helper test file

### No behavior change

Given:

- the tests are reorganized

Then:

- test behavior and coverage remain unchanged except for updates needed by the reorganization itself

### Maintainability

Given:

- a developer needs to update one command and its tests

Then:

- they can find the relevant tests without scanning one large catch-all file

## Decisions

1. Keep shared test helpers in a common helper file.
2. Split command tests into separate files aligned with the command source files where practical.
