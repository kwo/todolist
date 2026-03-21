# Planning

This project is split into planning and user story documents. These are not all strict sequential phases, and some user stories are independent of each other.

## Documents

- [Plan: MVP basic task management](mvp.md)
- [User Story: task metadata](user-story-metadata.md)
- [User Story: task filtering](user-story-filtering.md)
- [User Story: machine-readable JSON output](user-story-json.md)
- [User Story: task directory selection](user-story-directory-selection.md)
- [User Story: internal CLI parsing and inferred command values](user-story-cli-parser.md)
- [User Story: task directory configuration](user-story-directory-config.md)
- [User Story: parent-child task grouping](user-story-parent-child.md)
- [User Story: task dependencies](user-story-dep.md)
- [User Story: `init` command](user-story-init.md)

## Recommended implementation order

1. [Plan: MVP basic task management](mvp.md)
2. [User Story: task metadata](user-story-metadata.md)
3. [User Story: task filtering](user-story-filtering.md)
4. [User Story: task directory selection](user-story-directory-selection.md)
5. [User Story: internal CLI parsing and inferred command values](user-story-cli-parser.md)
6. [User Story: task directory configuration](user-story-directory-config.md)
7. [User Story: `init` command](user-story-init.md)
8. [User Story: machine-readable JSON output](user-story-json.md)
9. [User Story: parent-child task grouping](user-story-parent-child.md)
10. [User Story: task dependencies](user-story-dep.md)

## Dependency notes

- task metadata depends on the MVP
- task filtering depends on task metadata
- task directory selection depends on the MVP
- task directory configuration depends on the MVP and is best paired with task directory selection
- internal CLI parsing depends on the MVP and may require updating examples and command shapes in other CLI-oriented stories
- `init` depends on task directory configuration and is best paired with task directory selection
- machine-readable JSON output depends on the MVP and is otherwise mostly independent
- parent-child task grouping depends on the MVP
- task dependencies depend on the MVP and are best implemented after task metadata, since `ready` depends on dependencies having `status: done`
