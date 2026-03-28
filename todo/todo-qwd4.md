---
id: todo-qwd4
title: support Homebrew and Linuxbrew tap installation from release artifacts
status: todo
priority: 5
parents:
    - todo-96ke
depends:
    - todo-sscg
createdAt: "2026-03-28T18:50:55Z"
lastModified: "2026-03-28T18:50:55Z"
---

# User Story: support Homebrew and Linuxbrew installation from release artifacts

Publish release assets in a layout that works for both macOS Homebrew and Linuxbrew formulas in a personal tap.

## Acceptance criteria

- release asset naming is stable and formula-friendly
- tap formula can select macOS and Linux archives by CPU architecture
- install path works for Homebrew and Linuxbrew users from the same tap
- usage or release docs explain how to update the tap formula
