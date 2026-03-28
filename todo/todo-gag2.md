---
id: todo-gag2
title: add todolist version command using Go build metadata
status: todo
priority: 5
parents:
    - todo-96ke
createdAt: "2026-03-28T18:50:43Z"
lastModified: "2026-03-28T19:01:21Z"
---

# User Story: add version command using Go build metadata

Add a `todolist version` command that reports release version and embedded VCS metadata for released and local builds.

## Technical direction

- use `-ldflags` to inject the release semantic version at build time
- use Go build info / VCS stamping for commit and dirty state where available
- prefer standard Go build metadata over custom commit/date injection when possible

## Acceptance criteria

- `todolist version` prints the application semantic version
- release builds can report the tag-derived version value
- local builds have a sensible fallback version such as `dev`
- when VCS metadata is available, output includes commit revision and whether the tree was dirty
- command output remains useful when VCS metadata is unavailable
