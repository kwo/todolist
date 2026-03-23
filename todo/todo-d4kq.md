---
id: todo-d4kq
title: 'MVP: Basic todo management'
status: done
createdAt: "2026-03-23T18:50:56Z"
lastModified: "2026-03-23T18:50:56Z"
---

# MVP: Basic todo management

The MVP focuses on the smallest useful version of Todo List: basic todo CRUD backed by Markdown files.

## CLI shape

For the MVP, the CLI should optimize for fast human typing:

- use **arguments** for the main thing a command acts on
- use **command-first parsing**
- use **explicit `key=value` notation** only when disambiguation is needed
- use **stdin** for Markdown description input instead of a dedicated description flag

That means:

- todo IDs should be positional arguments for `view`, `update`, and `delete`
- the todo title should be a positional argument for `add`
- `update` should accept an optional replacement title as a command value
- `update` may also accept `title=<title>` for explicit assignment
- the todo description should be read from stdin when stdin is piped

## Commands

### `add`

Create a todo.

```bash
todolist add <title>
```

Arguments:

- `<title>` — required todo title

Flags:

- none in the MVP

Description input:

- if stdin is piped, the CLI should read the todo description from stdin
- if stdin is not piped, the todo should be created with an empty description

Output:

- on success, print the generated todo ID to stdout
- on error, print a message to stderr and exit with code 1

Examples:

```bash
todolist add "Buy groceries"
printf 'Need milk, eggs, and bread.\n' | todolist add "Buy groceries"
```

### `list`

List todos.

```bash
todolist list
```

Arguments:

- none

Flags:

- none in the MVP

Output:

- one todo per line, tab-separated: `<id>\t<title>`
- if there are no todos, print nothing

Example output:

```
todo-7k9m	Buy groceries
todo-2w8x	Write proposal
```

### `view`

Show a single todo.

```bash
todolist view <todo>
```

Arguments:

- `<todo>` — required todo ID

Flags:

- none in the MVP

Output:

- print the raw todo file as-is, including the YAML front matter and Markdown description
- do not render Markdown
- if the todo does not exist, print an error to stderr and exit with code 1

Example output:

```
---
id: todo-7k9m
title: Buy groceries
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

### `update`

Update a todo.

```bash
todolist update <todo> [<title>]
todolist update <todo> [title=<title>]
```

Arguments:

- `<todo>` — required todo ID
- `<title>` — optional replacement title

Flags:

- none in the MVP

Description input:

- if stdin is piped, the CLI should replace the todo description with stdin contents
- if stdin is not piped, the description should remain unchanged

Behavior:

- `id` is not updatable
- if a command value is provided and it is not otherwise recognized, it is assigned as the replacement title
- `title=<title>` may be used for explicit assignment
- the command should require at least one change, either a replacement title or piped stdin
- a successful update should automatically set `lastModified` to the current time

Output:

- on success, print nothing
- if the todo does not exist, print an error to stderr and exit with code 1
- if no changes are provided, print an error to stderr and exit with code 1

Examples:

```bash
todolist update todo-7k9m "Buy groceries and snacks"
todolist update todo-7k9m title="Buy groceries and snacks"
printf 'Need milk, eggs, bread, and chips.\n' | todolist update todo-7k9m
printf 'Need milk, eggs, bread, and chips.\n' | todolist update todo-7k9m title="Buy groceries and snacks"
```

### `delete`

Delete a todo.

```bash
todolist delete <todo>
```

Arguments:

- `<todo>` — required todo ID

Flags:

- none in the MVP

Behavior:

- delete the todo file immediately without prompting
- confirmation prompting and forced deletion behavior are deferred to the parent-child todo grouping story, where deleting a parent todo with children introduces cases that need protection

Output:

- on success, print nothing
- if the todo does not exist, print an error to stderr and exit with code 1

## Todo Format

Each todo is a Markdown file with YAML front matter and an optional Markdown description.

Example:

```md
---
id: todo-7k9m
title: Buy groceries
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
```

Front matter fields in the MVP:

- `id` — unique todo identifier; not updatable
- `title`
- `createdAt` — set on creation to the current time; uses RFC 3339 in UTC
- `lastModified` — set on creation to the current time; updated automatically on every successful `update`; uses RFC 3339 in UTC

On creation, `createdAt` and `lastModified` should both be set to the same current time.

The Markdown description is also part of the todo and is updatable.

The MVP does not write `status` or `priority` fields to todo files. Those are introduced in the todo metadata story.

## Error Handling

- nonexistent todo ID → error message to stderr, exit code 1
- missing or inaccessible todo directory → error message to stderr, exit code 1
- `update` with no replacement title and no description input provided → error message to stderr, exit code 1
- all other errors → message to stderr, exit code 1
- success → exit code 0

## Deferred to later stories

The following are explicitly **not** part of the MVP:

- `status` and `priority` metadata — see `user-story-metadata.md`
- filtering in `list` — see `user-story-filtering.md`
- `--json` global option — see `user-story-json.md`
- `-d, --directory` global option — see `user-story-directory-selection.md`
- `TODOLIST_DIRECTORY` environment variable — see `user-story-directory-selection.md`
- `.todos` config file and configurable ID prefix — see `user-story-directory-config.md`
- confirmation prompting and forced deletion controls on `delete` — see `user-story-parent-child.md`
- parent-child relationships — see `user-story-parent-child.md`
- dependencies — see `user-story-dep.md`

In the MVP, the todo directory is always `./todo` and the ID prefix is always `todo-`.

## Scope

- Markdown todo files with YAML front matter
- one todo per file
- basic CRUD commands: `add`, `list`, `view`, `update`, `delete`
- auto-generated todo IDs with hardcoded `todo-` prefix
- file naming as `<id>.md`
- todo directory is `./todo`
- human-readable CLI output
- description input through stdin for `add` and `update`
- timestamps in RFC 3339 UTC
- exit code 0 on success, 1 on error

Shared specifications such as ID generation, file naming, and principles are documented in [README.md](../README.md).

## Open Issues

None currently.
 [README.md](../README.md).

## Open Issues

None currently.
Issues

None currently.
