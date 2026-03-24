## Project Workflow

### Definition of Done

Before committing any changes, run this checklist in order:

```bash
golangci-lint run --fix
go test ./...
go mod tidy
git status
git add <files>
git commit -m "..."
```

### Best Practices

- After adding and using a new dependency, always run `go mod tidy`
- Always run `golangci-lint run --fix` before committing and fix any remaining issues
- After modifying app functionality, always update `USAGE.md` so the usage guide stays in sync with the CLI behavior

## Using the todolist app

Always run `todolist usage` to incorporate the app's usage instructions into the current context before managing todo items.
