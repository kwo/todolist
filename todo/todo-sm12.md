---
id: todo-sm12
title: remove configurable directory and hardcode ./todo
status: done
priority: 5
createdAt: "2026-03-28T18:42:44Z"
lastModified: "2026-03-28T19:08:53Z"
---

# User Story: remove configurable todo directory

Hardcode the todo directory to `./todo` and remove the configurable directory surface.

## Scope

- remove support for configurable todo directories
- remove `-d` and `--directory` flags
- stop reading `TODOLIST_DIRECTORY`
- hardcode directory resolution to `./todo`
- update usage/docs/tests accordingly
