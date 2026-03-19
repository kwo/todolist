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

- Use `github.com/jessevdk/go-flags` for CLI parsing
- Keep the implementation aligned with `README.md` and `docs/phase-1-mvp.md`
- After adding and using a new dependency, always run `go mod tidy`
- Always run `golangci-lint run --fix` before committing and fix any remaining issues
