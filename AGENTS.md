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

## Using the todolist app

For instructions on how to use the app to manage todo items, read `USAGE.md`.

Prefer `todolist ... --json` when an agent needs structured output.
