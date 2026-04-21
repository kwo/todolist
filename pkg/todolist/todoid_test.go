package todolist_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/kwo/todolist/pkg/todolist"
)

func TestGenerateIDStartsAtZeroWhenLastIDMissing(t *testing.T) {
	t.Helper()

	exists := func(string) bool { return false }

	generatedID, err := todolist.GenerateID("", exists)
	if err != nil {
		t.Fatalf("generate id: %v", err)
	}

	if generatedID != "todo-0000" {
		t.Fatalf("expected first id %q, got %q", "todo-0000", generatedID)
	}
}

func TestGenerateIDIncrementsFromLastID(t *testing.T) {
	t.Helper()

	exists := func(string) bool { return false }

	generatedID, err := todolist.GenerateID("todo-0000", exists)
	if err != nil {
		t.Fatalf("generate id: %v", err)
	}

	if generatedID != "todo-0001" {
		t.Fatalf("expected incremented id %q, got %q", "todo-0001", generatedID)
	}
}

func TestGenerateIDRetriesOnCollision(t *testing.T) {
	t.Helper()

	exists := func(id string) bool {
		return id == "todo-0001" || id == "todo-0002"
	}

	generatedID, err := todolist.GenerateID("todo-0000", exists)
	if err != nil {
		t.Fatalf("generate id: %v", err)
	}

	if generatedID != "todo-0003" {
		t.Fatalf("expected collision retry id %q, got %q", "todo-0003", generatedID)
	}
}

func TestGenerateIDWithPrefixUsesConfiguredPrefix(t *testing.T) {
	t.Helper()

	exists := func(string) bool { return false }

	generatedID, err := todolist.GenerateIDWithPrefix("", "work-", exists)
	if err != nil {
		t.Fatalf("generate id: %v", err)
	}

	matched, err := regexp.MatchString(`^work-[0-9abcdefghjkmnpqrstvwxyz]{4}$`, generatedID)
	if err != nil {
		t.Fatalf("match id format: %v", err)
	}

	if !matched {
		t.Fatalf("expected id to match work-xxxx format, got %q", generatedID)
	}
}

func TestGenerateIDReturnsSpaceExhaustedError(t *testing.T) {
	t.Helper()

	exists := func(string) bool { return false }

	_, err := todolist.GenerateID("todo-zzzz", exists)
	if !errors.Is(err, todolist.ErrIDSpaceExhausted) {
		t.Fatalf("expected ErrIDSpaceExhausted, got %v", err)
	}
}

func TestGenerateIDRejectsInvalidLastID(t *testing.T) {
	t.Helper()

	exists := func(string) bool { return false }

	_, err := todolist.GenerateID("work-0000", exists)
	if !errors.Is(err, todolist.ErrInvalidLastID) {
		t.Fatalf("expected ErrInvalidLastID for prefix mismatch, got %v", err)
	}

	_, err = todolist.GenerateID("todo-00i0", exists)
	if !errors.Is(err, todolist.ErrInvalidLastID) {
		t.Fatalf("expected ErrInvalidLastID for unsupported alphabet character, got %v", err)
	}
}
