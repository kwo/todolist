package cli_test

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestUpdateWithFlags(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	id := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Buy groceries"})

	exitCode := app.Run([]string{"update", id, "--title", "done"})
	if exitCode != 0 {
		t.Fatalf("expected update to succeed, got %d: %s", exitCode, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", id})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	viewed := stdout.String()
	if !strings.Contains(viewed, "title: done") {
		t.Fatalf("expected updated literal title, got %q", viewed)
	}
}

func TestUpdateRejectsInvalidPriority(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	id := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Buy groceries"})

	exitCode := app.Run([]string{"update", id, "--priority", "6"})
	if exitCode != 1 {
		t.Fatalf("expected update to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "invalid priority 6: must be between 1 and 5") {
		t.Fatalf("expected invalid priority error, got %q", stderr.String())
	}
}

func TestUpdateRequiresAtLeastOneChange(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	id := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Buy groceries"})

	exitCode := app.Run([]string{"update", id})
	if exitCode != 1 {
		t.Fatalf("expected update without changes to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "update requires a title, status, priority, parent, or stdin description input") {
		t.Fatalf("expected update error, got %q", stderr.String())
	}
}

func TestUpdateParentOperationsAndDeleteCleanup(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	parentOne := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Parent one"})
	parentTwo := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Parent two"})
	childID := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Child", "--parent", parentOne, "--parent", parentTwo})

	exitCode := app.Run([]string{"list"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	if !strings.Contains(stdout.String(), childID+"\t5\ttodo\tChild\t"+parentOne+",...") {
		t.Fatalf("expected list parents column, got %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"update", "--json", childID, "--parent", parentOne + "!"})
	if exitCode != 0 {
		t.Fatalf("expected remove parent update to succeed, got %d: %s", exitCode, stderr.String())
	}

	var updated jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &updated); err != nil {
		t.Fatalf("unmarshal update json: %v; output=%q", err, stdout.String())
	}

	if len(updated.Parents) != 1 || updated.Parents[0] != parentTwo {
		t.Fatalf("expected remaining parent %q, got %+v", parentTwo, updated.Parents)
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"update", childID, "--parent", parentOne + "!"})
	if exitCode != 1 {
		t.Fatalf("expected removing non-assigned parent to fail, got %d", exitCode)
	}

	if !strings.Contains(stderr.String(), "is not currently assigned") {
		t.Fatalf("expected non-assigned parent error, got %q", stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"delete", parentTwo})
	if exitCode != 0 {
		t.Fatalf("expected delete parent to succeed, got %d: %s", exitCode, stderr.String())
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", "--json", childID})
	if exitCode != 0 {
		t.Fatalf("expected json view to succeed, got %d: %s", exitCode, stderr.String())
	}

	var child jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &child); err != nil {
		t.Fatalf("unmarshal child json: %v; output=%q", err, stdout.String())
	}

	if len(child.Parents) != 0 {
		t.Fatalf("expected parent cleanup on delete, got %+v", child.Parents)
	}
}

func TestUpdateMissingTodoFailsBeforeNoChangeValidation(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"update", "todo-7k9m"})
	if exitCode != 1 {
		t.Fatalf("expected update to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `read todo "todo-7k9m"`) {
		t.Fatalf("expected missing todo error, got %q", stderr.String())
	}
}
