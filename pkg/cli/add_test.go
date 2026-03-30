package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJSONOmitsEmptyDescription(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "--json", "--title", "Buy groceries"})
	if exitCode != 0 {
		t.Fatalf("expected json add to succeed, got %d: %s", exitCode, stderr.String())
	}

	var added map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &added); err != nil {
		t.Fatalf("unmarshal add json: %v; output=%q", err, stdout.String())
	}

	if _, ok := added["description"]; ok {
		t.Fatalf("expected empty description to be omitted, got %+v", added)
	}
}

func TestAddDefaultsStatusAndPriority(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "--title", "Buy groceries"})
	if exitCode != 0 {
		t.Fatalf("expected add to succeed, got %d: %s", exitCode, stderr.String())
	}

	id := strings.TrimSpace(stdout.String())
	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", id})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	viewed := stdout.String()
	if !strings.Contains(viewed, "status: todo") {
		t.Fatalf("expected default status in view output, got %q", viewed)
	}

	if !strings.Contains(viewed, "priority: 5") {
		t.Fatalf("expected default priority in view output, got %q", viewed)
	}
}

func TestAddPositionalTitle(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "Buy groceries"})
	if exitCode != 0 {
		t.Fatalf("expected add to succeed, got %d: %s", exitCode, stderr.String())
	}

	id := strings.TrimSpace(stdout.String())
	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", id})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	if !strings.Contains(stdout.String(), "title: Buy groceries") {
		t.Fatalf("expected positional title, got %q", stdout.String())
	}
}

func TestAddUsesConfiguredPrefixFromTodoDirectory(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	if err := os.WriteFile(filepath.Join(app.TodoDir, ".todos"), []byte("prefix=work-\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	exitCode := app.Run([]string{"add", "--title", "Buy groceries"})
	if exitCode != 0 {
		t.Fatalf("expected add to succeed, got %d: %s", exitCode, stderr.String())
	}

	id := strings.TrimSpace(stdout.String())
	if !strings.HasPrefix(id, "work-") {
		t.Fatalf("expected configured prefix, got %q", id)
	}
}

func TestAddTitleFlagHandlesStatusLikeValues(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "--title", "done"})
	if exitCode != 0 {
		t.Fatalf("expected add to succeed, got %d: %s", exitCode, stderr.String())
	}

	id := strings.TrimSpace(stdout.String())
	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", id})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	viewed := stdout.String()
	if !strings.Contains(viewed, "title: done") {
		t.Fatalf("expected literal title in view output, got %q", viewed)
	}

	if !strings.Contains(viewed, "status: todo") {
		t.Fatalf("expected default status in view output, got %q", viewed)
	}
}

func TestAddTitleFlagHandlesPriorityLikeValues(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "--title", "2"})
	if exitCode != 0 {
		t.Fatalf("expected add to succeed, got %d: %s", exitCode, stderr.String())
	}

	id := strings.TrimSpace(stdout.String())
	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", id})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	viewed := stdout.String()
	if !strings.Contains(viewed, "title: \"2\"") {
		t.Fatalf("expected literal title in view output, got %q", viewed)
	}

	if !strings.Contains(viewed, "priority: 5") {
		t.Fatalf("expected default priority in view output, got %q", viewed)
	}
}

func TestAddRejectsInvalidStatus(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "--title", "Buy groceries", "--status", "active"})
	if exitCode != 1 {
		t.Fatalf("expected add to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `invalid status "active": must be one of todo, wip, done`) {
		t.Fatalf("expected invalid status error, got %q", stderr.String())
	}
}

func TestAddRejectsExtraPositionalArgs(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "Buy groceries", "snacks"})
	if exitCode != 1 {
		t.Fatalf("expected add to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "extra") {
		t.Fatalf("expected extra argument error, got %q", stderr.String())
	}
}

func TestAddRejectsBothPositionalAndFlagTitle(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "--title", "Buy groceries", "also groceries"})
	if exitCode != 1 {
		t.Fatalf("expected add to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "cannot use both") {
		t.Fatalf("expected conflict error, got %q", stderr.String())
	}
}

func TestAddAndViewParents(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	parentID := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Parent todo"})

	exitCode := app.Run([]string{"add", "--json", "--title", "Child todo", "--parent", parentID})
	if exitCode != 0 {
		t.Fatalf("expected add with parent to succeed, got %d: %s", exitCode, stderr.String())
	}

	var child jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &child); err != nil {
		t.Fatalf("unmarshal add json: %v; output=%q", err, stdout.String())
	}

	if len(child.Parents) != 1 || child.Parents[0] != parentID {
		t.Fatalf("expected child parents to include %q, got %+v", parentID, child.Parents)
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", "--json", parentID})
	if exitCode != 0 {
		t.Fatalf("expected parent json view to succeed, got %d: %s", exitCode, stderr.String())
	}

	var parent jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &parent); err != nil {
		t.Fatalf("unmarshal parent json: %v; output=%q", err, stdout.String())
	}

	if len(parent.Depends) != 1 || parent.Depends[0] != child.ID {
		t.Fatalf("expected parent depends to include child %q, got %+v", child.ID, parent.Depends)
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", child.ID})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	viewed := stdout.String()
	if !strings.Contains(viewed, "parents:") || !strings.Contains(viewed, parentID) {
		t.Fatalf("expected stored parents front matter, got %q", viewed)
	}

	if !strings.Contains(viewed, "Parents:\n- "+parentID+" Parent todo") {
		t.Fatalf("expected human-friendly parents section, got %q", viewed)
	}
}

func TestAddDependsDeduplicatesAndComputesReady(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	dependencyID := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Dependency", "--status", "done"})

	exitCode := app.Run([]string{"add", "--json", "--title", "Blocked todo", "--depends", dependencyID, "--depends", dependencyID})
	if exitCode != 0 {
		t.Fatalf("expected add with dependencies to succeed, got %d: %s", exitCode, stderr.String())
	}

	var added jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &added); err != nil {
		t.Fatalf("unmarshal add json: %v; output=%q", err, stdout.String())
	}

	if len(added.Depends) != 1 || added.Depends[0] != dependencyID {
		t.Fatalf("expected deduplicated dependency %q, got %+v", dependencyID, added.Depends)
	}

	if !added.Ready {
		t.Fatalf("expected ready to be true when dependency is done, got %+v", added)
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", added.ID})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	if !strings.Contains(stdout.String(), "depends:") || !strings.Contains(stdout.String(), dependencyID) {
		t.Fatalf("expected stored depends front matter, got %q", stdout.String())
	}
}

func TestAddHelpPrintsCommandUsage(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "--help"})
	if exitCode != 0 {
		t.Fatalf("expected help to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}

	if !strings.Contains(stdout.String(), "todolist add") {
		t.Fatalf("expected add usage, got %q", stdout.String())
	}
}
