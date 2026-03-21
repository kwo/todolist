# User Story: internal CLI parsing and inferred command values

Add an internal command parser to todolist and remove the external `github.com/jessevdk/go-flags` dependency.

## User story

As a user,
I want todolist to treat the first argument as the command and infer command-specific values from the remaining arguments,
so that the CLI is simpler to type and the application does not depend on `go-flags`.

## Goal

This work adds:

- an internal parser owned by todolist
- a command-first CLI shape
- support for parsing global options after the command
- inferred assignment of command-specific values without per-command flags such as `--status` and `--priority`
- explicit `title=...`, `priority=...`, and `status=...` notation for disambiguation
- removal of `github.com/jessevdk/go-flags` from the project

## Command surface

The root CLI shape becomes:

```bash
todolist <command> [global options] [command values...]
```

This story intentionally changes the command grammar used by earlier planning documents. If implemented, command examples in `README.md` and related user stories should be updated to the new command-first form.

Examples:

```bash
todolist add "Buy groceries"
todolist add "Buy groceries" wip 2
todolist add title=done
todolist list done
todolist list +3
todolist update todo-7k9m "Buy groceries and snacks" done 1
todolist update todo-7k9m title=done
todolist view todo-7k9m
todolist delete todo-7k9m
```

## Parsing model

### Root parsing

Given a process invocation:

```bash
todolist <command> ...
```

Then:

1. the first argument after the program name selects the command
2. if the first argument is missing, todolist returns a usage error
3. if the first argument is unknown, todolist returns an unknown command error
4. the remaining arguments are parsed as either recognized global options or command-specific values

### Global options

Global options remain global in meaning, but are parsed after the command.

Examples of intended usage:

```bash
todolist list --json
todolist update -d ./work todo-7k9m done
todolist add --json "Buy groceries"
```

This story applies to current and planned global options such as:

- `-d, --directory`
- `--json`
- `-h, --help`

## Inference rules

The parser should infer command-specific values from the command schema and the shape of each value.

### Recognizable value types

The parser should recognize these kinds of values:

- todo ID
  - a value that matches the todo ID format expected by todolist
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

To make ambiguous inputs deterministic, todolist should apply these rules:

1. command-required positional values are assigned first
   - example: `update` always assigns the first command value to `<todo>`
2. after required fields are assigned, remaining values are parsed left to right
3. any value using explicit key notation is assigned directly
   - `title=<value>` assigns the title
   - `status=<value>` assigns the status
   - `priority=<value>` assigns the priority
4. if a remaining value matches an unassigned `status` field, it is assigned to `status`
   - example: `update todo-7k9m done` sets `status: done`
5. if a remaining value matches an unassigned `priority` field, it is assigned to `priority`
   - example: `update todo-7k9m 2` sets `priority: 2`
6. any remaining value that is not a recognized status or priority value is assigned to the next unassigned free-text field
   - example: `add "Buy groceries" done 2` assigns title, then status, then priority
   - example: `update todo-7k9m "Buy groceries and snacks" done` assigns title first, then status
7. if the user wants a literal title that would otherwise be recognized as a status or priority, the user must use explicit notation
   - example: `add title=done` creates a todo with title `done`
   - example: `add title=2` creates a todo with title `2`
   - example: `update todo-7k9m title=done` sets the title to the literal value `done`
8. if a value cannot be assigned uniquely after these rules are applied, the command returns an error

These rules keep typed shorthand convenient while using explicit key notation to handle ambiguous values.

### Description input

This story does not change description handling.

- `add` and `update` should continue to read the todo description from stdin when stdin is piped
- description should not gain a new positional syntax as part of this story

## Command expectations

### `add`

```bash
todolist add <title> [<status>] [<priority>]
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
todolist add "Buy groceries"
todolist add "Buy groceries" done
todolist add "Buy groceries" done 2
todolist add title=done
todolist add title=2
```

### `update`

```bash
todolist update <todo> [<title>] [<status>] [<priority>]
```

Behavior:

- the first command-specific value is the todo ID
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
todolist update todo-7k9m "Buy groceries and snacks"
todolist update todo-7k9m done
todolist update todo-7k9m 1
todolist update todo-7k9m "Buy groceries and snacks" done 1
todolist update todo-7k9m title=done
todolist update todo-7k9m title=2
todolist update todo-7k9m status=done priority=2
```

### `list`

```bash
todolist list [<status-filter>] [<priority-filter>]
```

Behavior:

- status filters and priority filters should be recognized by their existing value syntax
- explicit `status=<value>` and `priority=<value>` notation may also be used
- command-specific filter flags should no longer be required

Examples:

```bash
todolist list done
todolist list !done
todolist list 3
todolist list +3
todolist list status=done
todolist list priority=+3
```

### `view`

```bash
todolist view <todo>
```

Behavior:

- the todo ID remains a required command value

### `delete`

```bash
todolist delete <todo>
```

Behavior:

- the todo ID remains a required command value

## Acceptance criteria

### Internal parser

Given:

- todolist parses CLI arguments

Then:

- parsing is handled by todolist code instead of `github.com/jessevdk/go-flags`
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

- a user runs todolist

Then:

- the first argument after the program name is always interpreted as the command
- recognized global options are parsed from the remaining arguments

### Inferred typed values

Given:

- a command accepts typed values such as todo ID, status, priority, or filters

Then:

- the parser assigns those values based on the command schema and the value format
- command-specific flags are not required for those values

### Ambiguity handling

Given:

- a value could map to more than one field

Then:

- todolist applies the documented precedence rules in this story
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
2. It is not yet specified whether todolist should support a temporary backwards-compatible migration period for the existing flag-based syntax.
3. It is not yet specified whether explicit key notation should allow quoted values such as `title="Buy groceries"` in addition to normal shell quoting like `title=Buy groceries` within quotes.
4. If this story is implemented, related planning documents and examples in `README.md` should be updated to match the new grammar.
