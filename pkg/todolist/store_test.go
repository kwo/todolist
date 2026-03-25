package todolist_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kwo/todolist/pkg/todolist"
)

func TestStoreGetDefaultsMissingMetadataFields(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	raw := `---
id: todo-7k9m
title: Buy groceries
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---

Need milk, eggs, and bread.
`

	path := filepath.Join(dir, "todo-7k9m.md")
	if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
		t.Fatalf("write todo file: %v", err)
	}

	value, err := todolist.NewStore(dir).Get("todo-7k9m")
	if err != nil {
		t.Fatalf("get todo: %v", err)
	}

	if value.Status != todolist.DefaultStatus {
		t.Fatalf("expected default status %q, got %q", todolist.DefaultStatus, value.Status)
	}

	if value.Priority != todolist.DefaultPriority {
		t.Fatalf("expected default priority %d, got %d", todolist.DefaultPriority, value.Priority)
	}
}

func TestStoreCreateSerializesMetadataFields(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	parentPath := filepath.Join(dir, "todo-parent.md")
	parentRaw := `---
id: todo-parent
title: Parent
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---
`
	if err := os.WriteFile(parentPath, []byte(parentRaw), 0o600); err != nil {
		t.Fatalf("write parent todo: %v", err)
	}

	value := todolist.Todo{
		ID:           "todo-7k9m",
		Title:        "Buy groceries",
		Status:       "wip",
		Priority:     2,
		Parents:      []string{"todo-parent"},
		CreatedAt:    time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC),
		LastModified: time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC),
		Description:  "Need milk, eggs, and bread.\n",
	}

	if err := todolist.NewStore(dir).Create(value); err != nil {
		t.Fatalf("create todo: %v", err)
	}

	//nolint:gosec // Test reads a file created in a temporary directory.
	raw, err := os.ReadFile(filepath.Join(dir, "todo-7k9m.md"))
	if err != nil {
		t.Fatalf("read todo file: %v", err)
	}

	text := string(raw)
	if !containsAll(text, "status: wip", "priority: 2", "parents:", "- todo-parent") {
		t.Fatalf("expected serialized metadata fields, got %q", text)
	}
}

func TestStoreGetRejectsMissingParentReference(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	raw := `---
id: todo-7k9m
title: Buy groceries
parents:
  - todo-parent
createdAt: 2026-03-18T10:00:00Z
lastModified: 2026-03-18T10:00:00Z
---
`

	path := filepath.Join(dir, "todo-7k9m.md")
	if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
		t.Fatalf("write todo file: %v", err)
	}

	_, err := todolist.NewStore(dir).Get("todo-7k9m")
	if err == nil {
		t.Fatal("expected missing parent error")
	}

	if !strings.Contains(err.Error(), `parent todo "todo-parent" does not exist`) {
		t.Fatalf("expected missing parent error, got %v", err)
	}
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(value, part) {
			return false
		}
	}

	return true
}
