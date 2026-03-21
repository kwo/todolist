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
	value := todolist.Todo{
		ID:           "todo-7k9m",
		Title:        "Buy groceries",
		Status:       "wip",
		Priority:     2,
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
	if !containsAll(text, "status: wip", "priority: 2") {
		t.Fatalf("expected serialized metadata fields, got %q", text)
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
