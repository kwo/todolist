# User Story: `init` command

Add an `init` command to bootstrap a todolist directory for a user.

## User story

As a user,
I want to run a single command to initialize a todolist directory,
so that I do not need to manually create the todo directory or the `.todos` config file.

## Goal

The command should create the todo storage directory and initialize the `.todos` config file used by later user stories.

## Command

```bash
todolist init
```

Later work may extend this with global directory selection, for example:

```bash
todolist init -d ./work-todos
```

## Expected behavior

When `todolist init` is run:

1. create the todo directory if it does not already exist
2. create a `.todos` config file inside that directory if it does not already exist
3. write the default config contents:

```text
prefix=todo-
```

## Success cases

### Fresh initialization

Given:

- the target todo directory does not exist

When:

- the user runs `todolist init`

Then:

- the todo directory is created
- the `.todos` file is created in that directory
- the `.todos` file contains:

```text
prefix=todo-
```

### Directory exists but config does not

Given:

- the todo directory already exists
- the `.todos` file does not exist

When:

- the user runs `todolist init`

Then:

- the existing todo directory is left in place
- the `.todos` file is created
- the `.todos` file contains:

```text
prefix=todo-
```

### Already initialized

Given:

- the todo directory already exists
- the `.todos` file already exists

When:

- the user runs `todolist init`

Then:

- the command should succeed without overwriting the existing `.todos` file
- the command should be idempotent

## Error cases

- if the target path exists but is not a directory, return an error
- if the `.todos` path exists but is not a regular file, return an error
- if the directory or config file cannot be created, return an error

## Output

Human-readable output should confirm what happened. Example:

```text
initialized todo directory: ./todo
created config file: ./todo/.todos
```

If already initialized, output should make that clear. Example:

```text
todo directory already exists: ./todo
config file already exists: ./todo/.todos
```

## Notes

- This command is primarily useful once `.todos` configuration is supported.
- In the MVP, the default todo directory is `./todo`, so `init` naturally complements that workflow.
- The command should be safe to run multiple times.
