# Planning

This project is split into planning and user story documents. These are not all strict sequential phases, and some user stories are independent of each other.

## Documents

- [Plan: MVP basic task management](mvp.md)
- [User Story: task metadata](user-story-metadata.md)
- [User Story: task filtering](user-story-filtering.md)
- [User Story: machine-readable JSON output](user-story-json.md)
- [User Story: task directory selection](user-story-directory-selection.md)
- [User Story: task directory configuration](user-story-directory-config.md)
- [User Story: parent-child task grouping](user-story-parent-child.md)
- [User Story: task dependencies](user-story-dep.md)
- [User Story: `init` command](user-story-init.md)

## Recommended implementation order

1. [Plan: MVP basic task management](mvp.md)
2. [User Story: task metadata](user-story-metadata.md)
3. [User Story: task filtering](user-story-filtering.md)
4. [User Story: task directory selection](user-story-directory-selection.md)
5. [User Story: task directory configuration](user-story-directory-config.md)
6. [User Story: `init` command](user-story-init.md)
7. [User Story: machine-readable JSON output](user-story-json.md)
8. [User Story: parent-child task grouping](user-story-parent-child.md)
9. [User Story: task dependencies](user-story-dep.md)

## Dependency notes

- task metadata depends on the MVP
- task filtering depends on task metadata
- task directory selection depends on the MVP
- task directory configuration depends on the MVP and is best paired with task directory selection
- `init` depends on task directory configuration and is best paired with task directory selection
- machine-readable JSON output depends on the MVP and is otherwise mostly independent
- parent-child task grouping depends on the MVP
- task dependencies depend on the MVP and are best implemented after task metadata, since `ready` likely depends on `status: done`
