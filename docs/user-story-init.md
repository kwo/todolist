# User Story: `init` command

Add an `init` command to bootstrap a tasklist directory for a user.

## User story

As a user,
I want to run a single command to initialize a tasklist directory,
so that I do not need to manually create the tasks directory or the `.tasks` config file.

## Goal

The command should create the task storage directory and initialize the `.tasks` config file used by later user stories.

## Command

```bash
tasklist init
```

Later work may extend this with global directory selection, for example:

```bash
tasklist init -d ./work-tasks
```

## Expected behavior

When `tasklist init` is run:

1. create the task directory if it does not already exist
2. create a `.tasks` config file inside that directory if it does not already exist
3. write the default config contents:

```text
prefix=task-
```

## Success cases

### Fresh initialization

Given:

- the target task directory does not exist

When:

- the user runs `tasklist init`

Then:

- the task directory is created
- the `.tasks` file is created in that directory
- the `.tasks` file contains:

```text
prefix=task-
```

### Directory exists but config does not

Given:

- the task directory already exists
- the `.tasks` file does not exist

When:

- the user runs `tasklist init`

Then:

- the existing task directory is left in place
- the `.tasks` file is created
- the `.tasks` file contains:

```text
prefix=task-
```

### Already initialized

Given:

- the task directory already exists
- the `.tasks` file already exists

When:

- the user runs `tasklist init`

Then:

- the command should succeed without overwriting the existing `.tasks` file
- the command should be idempotent

## Error cases

- if the target path exists but is not a directory, return an error
- if the `.tasks` path exists but is not a regular file, return an error
- if the directory or config file cannot be created, return an error

## Output

Human-readable output should confirm what happened. Example:

```text
initialized task directory: ./tasks
created config file: ./tasks/.tasks
```

If already initialized, output should make that clear. Example:

```text
task directory already exists: ./tasks
config file already exists: ./tasks/.tasks
```

## Notes

- This command is primarily useful once `.tasks` configuration is supported.
- In the MVP, the default task directory is `./tasks`, so `init` naturally complements that workflow.
- The command should be safe to run multiple times.
