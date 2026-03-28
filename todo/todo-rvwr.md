---
id: todo-rvwr
title: build cross-platform release artifacts in GitHub Actions
status: todo
priority: 5
parents:
    - todo-96ke
depends:
    - todo-gag2
createdAt: "2026-03-28T18:50:47Z"
lastModified: "2026-03-28T18:50:47Z"
---

# User Story: build cross-platform release artifacts in GitHub Actions

Add a GitHub Actions workflow that builds release binaries for macOS and Linux on amd64 and arm64 when a version tag is pushed.

## Acceptance criteria

- tag push triggers the workflow
- builds target darwin/amd64, darwin/arm64, linux/amd64, linux/arm64
- packages each binary into a release archive
- emits sha256 checksums for all archives
