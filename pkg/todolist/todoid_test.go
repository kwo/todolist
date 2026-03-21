package todolist_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/kwo/todolist/pkg/todolist"
)

func TestGenerateIDIsDeterministicAndWellFormed(t *testing.T) {
	t.Helper()

	value := todolist.Todo{
		Title:       "Buy groceries",
		Description: "Need milk",
		CreatedAt:   time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC),
	}
	exists := func(string) bool { return false }

	first := todolist.GenerateID(value, exists)
	second := todolist.GenerateID(value, exists)

	if first != second {
		t.Fatalf("expected deterministic id, got %q and %q", first, second)
	}

	matched, err := regexp.MatchString(`^todo-[0-9a-z]{4}$`, first)
	if err != nil {
		t.Fatalf("match id format: %v", err)
	}

	if !matched {
		t.Fatalf("expected id to match todo-xxxx format, got %q", first)
	}
}

func TestGenerateIDIgnoresSubsecondCreatedAtPrecision(t *testing.T) {
	t.Helper()

	firstTodo := todolist.Todo{
		Title:       "Buy groceries",
		Description: "Need milk",
		CreatedAt:   time.Date(2026, time.March, 18, 10, 0, 0, 123000000, time.UTC),
	}
	secondTodo := todolist.Todo{
		Title:       "Buy groceries",
		Description: "Need milk",
		CreatedAt:   time.Date(2026, time.March, 18, 10, 0, 0, 987000000, time.UTC),
	}
	exists := func(string) bool { return false }

	firstID := todolist.GenerateID(firstTodo, exists)
	secondID := todolist.GenerateID(secondTodo, exists)

	if firstID != secondID {
		t.Fatalf("expected same id for timestamps within the same second, got %q and %q", firstID, secondID)
	}
}

func TestGenerateIDRetriesOnCollision(t *testing.T) {
	t.Helper()

	value := todolist.Todo{
		Title:       "Buy groceries",
		Description: "Need milk",
		CreatedAt:   time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC),
	}
	attempts := 0
	exists := func(string) bool {
		attempts++

		return attempts == 1
	}

	generatedID := todolist.GenerateID(value, exists)

	if attempts < 2 {
		t.Fatalf("expected at least two attempts, got %d", attempts)
	}

	if generatedID == "" {
		t.Fatal("expected generated id")
	}
}
