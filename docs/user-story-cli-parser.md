# User Story: internal CLI parsing and inferred command values

Add an internal command parser to tasklist and remove the external `github.com/jessevdk/go-flags` dependency.

## User story

As a user,
I want tasklist to treat the first argument as the command and infer command-specific values from the remaining arguments,
so that the CLI is simpler to type and the application does not depend on `go-flags`.

## Goal

This work adds:

- an internal parser owned by tasklist
- a command-first CLI shape
- support for parsing global options after the command
- inferred assignment of command-specific values without per-command flags such as `--status` and `--priority`
- explicit `title=...`, `priority=...`, and `status=...` notation for disambiguation
- removal of `github.com/jessevdk/go-flags` from the project

## Command surface

The root CLI shape becomes:

```bash
tasklist <command> [global options] [command values...]
```

This story intentionally changes the command grammar used by earlier planning documents. If implemented, command examples in `README.md` and related user stories should be updated to the new command-first form.

Examples:

```bash
tasklist add "Buy groceries"
tasklist add "Buy groceries" wip 2
tasklist add title=done
tasklist list done
tasklist list +3
tasklist update task-7k9m "Buy groceries and snacks" done 1
tasklist update task-7k9m title=done
tasklist view task-7k9m
tasklist delete task-7k9m
```

## Parsing model

### Root parsing

Given a process invocation:

```bash
tasklist <command> ...
```

Then:

1. the first argument after the program name selects the command
2. if the first argument is missing, tasklist returns a usage error
3. if the first argument is unknown, tasklist returns an unknown command error
4. the remaining arguments are parsed as either recognized global options or command-specific values

### Global options

Global options remain global in meaning, but are parsed after the command.

Examples of intended usage:

```bash
tasklist list --json
tasklist update -d ./work task-7k9m done
tasklist add --json "Buy groceries"
```

This story applies to current and planned global options such as:

- `-d, --directory`
- `--json`
- `-h, --help`

## Inference rules

The parser should infer command-specific values from the command schema and the shape of each value.

### Recognizable value types

The parser should recognize these kinds of values:

- task ID
  - a value that matches the task ID format expected by tasklist
- status
  - one of `todo`, `wip`, or `done`
- priority
  - an integer from `1` through `5`
- priority filter
  - values already defined by the filtering story, such as `3`, `.3`, `+3`, or `-3`
- status filter
  - values already defined by the filtering story, such as `done` or `!done`
- free-text value
  - any remaining value that does not match a recognized typed value and can be assigned to a free-text field such as `title`

### Assignment behavior

Given a command definition,
Then:

- required typed values are assigned first
- optional typed values are assigned by validation against supported value sets
- at most one free-text value is assigned to each free-text field unless the command explicitly supports more
- if a value is not a recognized status or priority value, it is assumed to be a free-text value such as `title`
- explicit `title=...`, `status=...`, and `priority=...` notation may be used to disambiguate intent
- if a value still cannot be assigned under these rules, the command returns an error instead of guessing silently

### Precedence and explicit-value rules

To make ambiguous inputs deterministic, tasklist should apply these rules:

1. command-required positional values are assigned first
   - example: `update` always assigns the first command value to `<task>`
2. after required fields are assigned, remaining values are parsed left to right
3. any value using explicit key notation is assigned directly
   - `title=<value>` assigns the title
   - `status=<value>` assigns the status
   - `priority=<value>` assigns the priority
4. if a remaining value matches an unassigned `status` field, it is assigned to `status`
   - example: `update task-7k9m done` sets `status: done`
5. if a remaining value matches an unassigned `priority` field, it is assigned to `priority`
   - example: `update task-7k9m 2` sets `priority: 2`
6. any remaining value that is not a recognized status or priority value is assigned to the next unassigned free-text field
   - example: `add "Buy groceries" done 2` assigns title, then status, then priority
   - example: `update task-7k9m "Buy groceries and snacks" done` assigns title first, then status
7. if the user wants a literal title that would otherwise be recognized as a status or priority, the user must use explicit notation
   - example: `add title=done` creates a task with title `done`
   - example: `add title=2` creates a task with title `2`
   - example: `update task-7k9m title=done` sets the title to the literal value `done`
8. if a value cannot be assigned uniquely after these rules are applied, the command returns an error

These rules keep typed shorthand convenient while using explicit key notation to handle ambiguous values.

### Description input

This story does not change description handling.

- `add` and `update` should continue to read the task description from stdin when stdin is piped
- description should not gain a new positional syntax as part of this story

## Command expectations

### `add`

```bash
tasklist add <title> [<status>] [<priority>]
```

Behavior:

- the first command-specific value is parsed left to right using the precedence rules in this story
- a value matching `todo`, `wip`, or `done` is assigned to `status` if `status` is still unassigned
- a value matching `1` through `5` is assigned to `priority` if `priority` is still unassigned
- any other value is assigned to `title` if `title` is still unassigned
- to use a literal title that looks like a status or priority, the user must write `title=<value>`
- if `status` is omitted, it defaults to `todo`
- if `priority` is omitted, it defaults to `5`

Examples:

```bash
tasklist add "Buy groceries"
tasklist add "Buy groceries" done
tasklist add "Buy groceries" done 2
tasklist add title=done
tasklist add title=2
```

### `update`

```bash
tasklist update <task> [<title>] [<status>] [<priority>]
```

Behavior:

- the first command-specific value is the task ID
- for the remaining values, inference is applied left to right
- values matching `todo`, `wip`, or `done` are assigned to `status` if `status` is still unassigned
- values matching `1` through `5` are assigned to `priority` if `priority` is still unassigned
- other values are assigned to `title` if `title` is still unassigned
- to use a literal title that looks like a status or priority, the user must write `title=<value>`
- explicit `status=<value>` and `priority=<value>` notation may also be used
- the command must still require at least one effective change, or piped stdin description input
- `lastModified` should still update automatically on success

Examples:

```bash
tasklist update task-7k9m "Buy groceries and snacks"
tasklist update task-7k9m done
tasklist update task-7k9m 1
tasklist update task-7k9m "Buy groceries and snacks" done 1
tasklist update task-7k9m title=done
tasklist update task-7k9m title=2
tasklist update task-7k9m status=done priority=2
```

### `list`

```bash
tasklist list [<status-filter>] [<priority-filter>]
```

Behavior:

- status filters and priority filters should be recognized by their existing value syntax
- explicit `status=<value>` and `priority=<value>` notation may also be used
- command-specific filter flags should no longer be required

Examples:

```bash
tasklist list done
tasklist list !done
tasklist list 3
tasklist list +3
tasklist list status=done
tasklist list priority=+3
```

### `view`

```bash
tasklist view <task>
```

Behavior:

- the task ID remains a required command value

### `delete`

```bash
tasklist delete <task>
```

Behavior:

- the task ID remains a required command value

## Acceptance criteria

### Internal parser

Given:

- tasklist parses CLI arguments

Then:

- parsing is handled by tasklist code instead of `github.com/jessevdk/go-flags`
- command definitions and dispatch remain testable
- usage and validation errors remain human-readable

### Dependency removal

Given:

- the implementation is complete

Then:

- `go.mod` no longer requires `github.com/jessevdk/go-flags`
- `go.sum` no longer contains entries for `github.com/jessevdk/go-flags`

### Command-first invocation

Given:

- a user runs tasklist

Then:

- the first argument after the program name is always interpreted as the command
- recognized global options are parsed from the remaining arguments

### Inferred typed values

Given:

- a command accepts typed values such as task ID, status, priority, or filters

Then:

- the parser assigns those values based on the command schema and the value format
- command-specific flags are not required for those values

### Ambiguity handling

Given:

- a value could map to more than one field

Then:

- tasklist applies the documented precedence rules in this story
- recognized status and priority values are inferred by default
- values that are not recognized as status or priority are treated as free-text values such as `title`
- `title=<value>`, `status=<value>`, and `priority=<value>` allow the user to state intent explicitly
- if a value still cannot be assigned under those rules, the parser returns a clear error

## Scope

- internal CLI parsing implementation
- command-first grammar
- global option parsing after the command
- inferred command-specific value assignment
- removal of the `go-flags` dependency
- updates to tests covering the new parsing behavior

## Open issues

1. It is not yet specified whether recognized global options may appear anywhere after the command or only before the first command-specific value.
2. It is not yet specified whether tasklist should support a temporary backwards-compatible migration period for the existing flag-based syntax.
3. It is not yet specified whether explicit key notation should allow quoted values such as `title="Buy groceries"` in addition to normal shell quoting like `title=Buy groceries` within quotes.
4. If this story is implemented, related planning documents and examples in `README.md` should be updated to match the new grammar.
