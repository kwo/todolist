---
id: todo-sscg
title: publish GitHub releases from pushed tags
status: todo
priority: 5
parents:
    - todo-96ke
depends:
    - todo-rvwr
createdAt: "2026-03-28T18:50:51Z"
lastModified: "2026-03-28T18:50:51Z"
---

# User Story: publish GitHub releases from tags

Extend the release workflow so a pushed version tag creates a GitHub Release and uploads all release artifacts automatically.

## Acceptance criteria

- release workflow creates a GitHub Release for the pushed tag
- uploaded assets include all archives and checksums
- a maintainer can publish by pushing only a git tag
