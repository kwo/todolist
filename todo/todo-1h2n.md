---
id: todo-1h2n
title: Migrate CLI to spf13/pflag with explicit flags
status: done
priority: 2
createdAt: "2026-03-24T21:08:50Z"
lastModified: "2026-03-24T21:08:50Z"
---

Replace custom positional/inference-based argument parsing with spf13/pflag. Each command gets its own FlagSet with explicit flags (--title/-t, --status/-s, --priority/-p) plus shared global flags (--directory/-d, --json, --help/-h). Removes ambiguous value inference and key=value assignment syntax. Updates help text, tests, and USAGE.md.
