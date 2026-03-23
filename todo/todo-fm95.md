---
id: todo-fm95
title: Planning
status: done
createdAt: "2026-03-23T18:50:56Z"
lastModified: "2026-03-23T18:50:56Z"
---

# Planning

This project is split into planning and user story documents. These are not all strict sequential phases, and some user stories are independent of each other.

## Documents

- [Plan: MVP basic todo management](mvp.md)
- [User Story: todo metadata](user-story-metadata.md)
- [User Story: todo filtering](user-story-filtering.md)
- [User Story: machine-readable JSON output](user-story-json.md)
- [User Story: todo directory selection](user-story-directory-selection.md)
- [User Story: internal CLI parsing and inferred command values](user-story-cli-parser.md)
- [User Story: todo directory configuration](user-story-directory-config.md)
- [User Story: parent-child todo grouping](user-story-parent-child.md)
- [User Story: todo dependencies](user-story-dep.md)
- [User Story: `init` command](user-story-init.md)

## Recommended implementation order

1. [Plan: MVP basic todo management](mvp.md)
2. [User Story: todo metadata](user-story-metadata.md)
3. [User Story: todo filtering](user-story-filtering.md)
4. [User Story: todo directory selection](user-story-directory-selection.md)
5. [User Story: internal CLI parsing and inferred command values](user-story-cli-parser.md)
6. [User Story: todo directory configuration](user-story-directory-config.md)
7. [User Story: `init` command](user-story-init.md)
8. [User Story: machine-readable JSON output](user-story-json.md)
9. [User Story: parent-child todo grouping](user-story-parent-child.md)
10. [User Story: todo dependencies](user-story-dep.md)

## Dependency notes

- todo metadata depends on the MVP
- todo filtering depends on todo metadata
- todo directory selection depends on the MVP
- todo directory configuration depends on the MVP and is best paired with todo directory selection
- internal CLI parsing depends on the MVP and may require updating examples and command shapes in other CLI-oriented stories
- `init` depends on todo directory configuration and is best paired with todo directory selection
- machine-readable JSON output depends on the MVP and is otherwise mostly independent
- parent-child todo grouping depends on the MVP
- todo dependencies depend on the MVP and are best implemented after todo metadata, since `ready` depends on dependencies having `status: done`
