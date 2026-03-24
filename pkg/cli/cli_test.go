package cli_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kwo/todolist/pkg/cli"
)

type jsonTodo struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Status       string `json:"status"`
	Priority     int    `json:"priority"`
	CreatedAt    string `json:"createdAt"`
	LastModified string `json:"lastModified"`
	Description  string `json:"description"`
}

type jsonDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

//nolint:maintidx // End-to-end CLI flow assertions are intentionally kept together.
func TestAddListViewUpdateDeleteFlow(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, true, "Need milk, eggs, and bread.\n")

	exitCode := app.Run([]string{"add", "--title", "Buy groceries", "--status", "wip", "--priority", "2"})
	if exitCode != 0 {
		t.Fatalf("expected add to succeed, got %d: %s", exitCode, stderr.String())
	}

	id := strings.TrimSpace(stdout.String())
	if id == "" {
		t.Fatal("expected add to print a todo id")
	}

	stdout.Reset()
	stderr.Reset()
	app.StdinProvided = false
	app.Stdin = strings.NewReader("")

	exitCode = app.Run([]string{"list"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	expectedLine := id + "\t2\twip\tBuy groceries\n"
	if stdout.String() != expectedLine {
		t.Fatalf("expected list output %q, got %q", expectedLine, stdout.String())
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", id})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	viewed := stdout.String()
	if !strings.Contains(viewed, "title: Buy groceries") {
		t.Fatalf("expected view output to contain title, got %q", viewed)
	}

	if !strings.Contains(viewed, "status: wip") {
		t.Fatalf("expected view output to contain status, got %q", viewed)
	}

	if !strings.Contains(viewed, "priority: 2") {
		t.Fatalf("expected view output to contain priority, got %q", viewed)
	}

	if !strings.Contains(viewed, "Need milk, eggs, and bread.\n") {
		t.Fatalf("expected view output to contain description, got %q", viewed)
	}

	stdout.Reset()
	stderr.Reset()
	app.StdinProvided = true
	app.Stdin = strings.NewReader("Need milk, eggs, bread, and chips.\n")

	exitCode = app.Run([]string{"update", id, "--title", "Buy groceries and snacks", "--status", "done", "--priority", "1"})
	if exitCode != 0 {
		t.Fatalf("expected update to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected update to print nothing, got %q", stdout.String())
	}

	//nolint:gosec // Test reads a file created in a temporary directory by the CLI.
	rawTodo, err := os.ReadFile(filepath.Join(app.TodoDir, id+".md"))
	if err != nil {
		t.Fatalf("read updated todo: %v", err)
	}

	updated := string(rawTodo)
	if !strings.Contains(updated, "title: Buy groceries and snacks") {
		t.Fatalf("expected updated title, got %q", updated)
	}

	if !strings.Contains(updated, "status: done") {
		t.Fatalf("expected updated status, got %q", updated)
	}

	if !strings.Contains(updated, "priority: 1") {
		t.Fatalf("expected updated priority, got %q", updated)
	}

	if !strings.Contains(updated, "Need milk, eggs, bread, and chips.\n") {
		t.Fatalf("expected updated description, got %q", updated)
	}

	stdout.Reset()
	stderr.Reset()
	app.StdinProvided = false
	app.Stdin = strings.NewReader("")

	exitCode = app.Run([]string{"delete", id})
	if exitCode != 0 {
		t.Fatalf("expected delete to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected delete to print nothing, got %q", stdout.String())
	}

	entries, err := os.ReadDir(app.TodoDir)
	if err != nil {
		t.Fatalf("read todo dir: %v", err)
	}

	if len(entries) != 0 {
		t.Fatalf("expected todo dir to be empty, found %d entries", len(entries))
	}
}

//nolint:maintidx // End-to-end JSON flow assertions are intentionally kept together.
func TestJSONOutputForCoreCommands(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, true, "Need milk, eggs, and bread.\n")

	exitCode := app.Run([]string{"add", "--json", "--title", "Buy groceries", "--status", "wip", "--priority", "2"})
	if exitCode != 0 {
		t.Fatalf("expected json add to succeed, got %d: %s", exitCode, stderr.String())
	}

	var added jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &added); err != nil {
		t.Fatalf("unmarshal add json: %v; output=%q", err, stdout.String())
	}

	if added.ID == "" {
		t.Fatal("expected add json to include an id")
	}

	if added.Title != "Buy groceries" || added.Status != "wip" || added.Priority != 2 {
		t.Fatalf("unexpected add json: %+v", added)
	}

	if added.Description != "Need milk, eggs, and bread.\n" {
		t.Fatalf("expected add description in json, got %q", added.Description)
	}

	stdout.Reset()
	stderr.Reset()
	app.StdinProvided = false
	app.Stdin = strings.NewReader("")

	exitCode = app.Run([]string{"list", "--json"})
	if exitCode != 0 {
		t.Fatalf("expected json list to succeed, got %d: %s", exitCode, stderr.String())
	}

	var listed []jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &listed); err != nil {
		t.Fatalf("unmarshal list json: %v; output=%q", err, stdout.String())
	}

	if len(listed) != 1 || listed[0].ID != added.ID {
		t.Fatalf("unexpected list json: %+v", listed)
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", "--json", added.ID})
	if exitCode != 0 {
		t.Fatalf("expected json view to succeed, got %d: %s", exitCode, stderr.String())
	}

	var viewed jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &viewed); err != nil {
		t.Fatalf("unmarshal view json: %v; output=%q", err, stdout.String())
	}

	if viewed != added {
		t.Fatalf("expected view json to match add json, got %+v want %+v", viewed, added)
	}

	stdout.Reset()
	stderr.Reset()
	app.StdinProvided = true
	app.Stdin = strings.NewReader("Need milk, eggs, bread, and chips.\n")

	exitCode = app.Run([]string{"update", "--json", added.ID, "--title", "Buy groceries and snacks", "--status", "done", "--priority", "1"})
	if exitCode != 0 {
		t.Fatalf("expected json update to succeed, got %d: %s", exitCode, stderr.String())
	}

	var updatedTodo jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &updatedTodo); err != nil {
		t.Fatalf("unmarshal update json: %v; output=%q", err, stdout.String())
	}

	if updatedTodo.ID != added.ID || updatedTodo.Title != "Buy groceries and snacks" || updatedTodo.Status != "done" || updatedTodo.Priority != 1 {
		t.Fatalf("unexpected update json: %+v", updatedTodo)
	}

	if updatedTodo.Description != "Need milk, eggs, bread, and chips.\n" {
		t.Fatalf("expected updated description in json, got %q", updatedTodo.Description)
	}

	stdout.Reset()
	stderr.Reset()
	app.StdinProvided = false
	app.Stdin = strings.NewReader("")

	exitCode = app.Run([]string{"delete", "--json", added.ID})
	if exitCode != 0 {
		t.Fatalf("expected json delete to succeed, got %d: %s", exitCode, stderr.String())
	}

	var deleted jsonDeleteResult
	if err := json.Unmarshal(stdout.Bytes(), &deleted); err != nil {
		t.Fatalf("unmarshal delete json: %v; output=%q", err, stdout.String())
	}

	if deleted.ID != added.ID || !deleted.Deleted {
		t.Fatalf("unexpected delete json: %+v", deleted)
	}
}

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

	if !strings.Contains(stderr.String(), "update requires a title, status, priority, or stdin description input") {
		t.Fatalf("expected update error, got %q", stderr.String())
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

func TestListIncludesIDPriorityStatusAndTitleColumns(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	id := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Buy groceries", "--status", "wip", "--priority", "2"})

	exitCode := app.Run([]string{"list"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	expectedLine := id + "\t2\twip\tBuy groceries\n"
	if stdout.String() != expectedLine {
		t.Fatalf("expected list output %q, got %q", expectedLine, stdout.String())
	}
}

func TestListTruncatesLongTitlesWithEllipsis(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	longTitle := "Investigate how to reconcile customer billing exports across regions and vendors"
	id := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", longTitle, "--status", "todo", "--priority", "3"})

	exitCode := app.Run([]string{"list"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	fields := strings.Split(strings.TrimSuffix(stdout.String(), "\n"), "\t")
	if len(fields) != 4 {
		t.Fatalf("expected 4 list columns, got %d in %q", len(fields), stdout.String())
	}

	if fields[0] != id || fields[1] != "3" || fields[2] != "todo" {
		t.Fatalf("unexpected list columns %q", stdout.String())
	}

	if len(fields[3]) != 60 {
		t.Fatalf("expected truncated title length 60, got %d in %q", len(fields[3]), fields[3])
	}

	if !strings.HasSuffix(fields[3], "...") {
		t.Fatalf("expected truncated title to end with ellipsis, got %q", fields[3])
	}

	if strings.Contains(fields[3], "vendors") {
		t.Fatalf("expected truncated title not to include full title, got %q", fields[3])
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", "--json", id})
	if exitCode != 0 {
		t.Fatalf("expected json view to succeed, got %d: %s", exitCode, stderr.String())
	}

	var viewed jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &viewed); err != nil {
		t.Fatalf("unmarshal view json: %v; output=%q", err, stdout.String())
	}

	if viewed.Title != longTitle {
		t.Fatalf("expected full title in json, got %q", viewed.Title)
	}
}

func TestListExcludesDoneByDefault(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "first todo", "--status", "todo"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "second todo", "--status", "done"})

	exitCode := app.Run([]string{"list"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "first todo") {
		t.Fatalf("expected first todo to be included, got %q", output)
	}

	if strings.Contains(output, "second todo") {
		t.Fatalf("expected done todo to be excluded by default, got %q", output)
	}
}

func TestListFiltersByStatus(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "first todo", "--status", "todo"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "second todo", "--status", "done"})

	exitCode := app.Run([]string{"list", "--status", "done"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if strings.Contains(output, "first todo") {
		t.Fatalf("expected first todo to be filtered out, got %q", output)
	}

	if !strings.Contains(output, "second todo") {
		t.Fatalf("expected done todo to be included, got %q", output)
	}
}

func TestListExcludesStatusWithBangSuffix(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "first todo", "--status", "todo"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "second todo", "--status", "done"})

	exitCode := app.Run([]string{"list", "--status", "done!"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "first todo") {
		t.Fatalf("expected first todo to be included, got %q", output)
	}

	if strings.Contains(output, "second todo") {
		t.Fatalf("expected done todo to be excluded, got %q", output)
	}
}

func TestListFiltersByPriority(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "high priority", "--priority", "1"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "low priority", "--priority", "5"})

	exitCode := app.Run([]string{"list", "--priority", "1"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "high priority") {
		t.Fatalf("expected priority 1 todo to be included, got %q", output)
	}

	if strings.Contains(output, "low priority") {
		t.Fatalf("expected priority 5 todo to be filtered out, got %q", output)
	}
}

func TestListPriorityFilterStillExcludesDoneByDefault(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "active priority one", "--status", "todo", "--priority", "1"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "done priority one", "--status", "done", "--priority", "1"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "active priority five", "--status", "todo", "--priority", "5"})

	exitCode := app.Run([]string{"list", "--priority", "1"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "active priority one") {
		t.Fatalf("expected active priority 1 todo to be included, got %q", output)
	}

	if strings.Contains(output, "done priority one") {
		t.Fatalf("expected done priority 1 todo to be excluded by default, got %q", output)
	}

	if strings.Contains(output, "active priority five") {
		t.Fatalf("expected priority 5 todo to be filtered out, got %q", output)
	}
}

func TestListPriorityFilterStillExcludesDoneByDefaultInJSON(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "active priority one", "--status", "todo", "--priority", "1"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "done priority one", "--status", "done", "--priority", "1"})

	exitCode := app.Run([]string{"list", "--json", "--priority", "1"})
	if exitCode != 0 {
		t.Fatalf("expected json list to succeed, got %d: %s", exitCode, stderr.String())
	}

	var listed []jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &listed); err != nil {
		t.Fatalf("unmarshal list json: %v; output=%q", err, stdout.String())
	}

	if len(listed) != 1 || listed[0].Title != "active priority one" {
		t.Fatalf("unexpected json list output: %+v", listed)
	}
}

func TestListFiltersPriorityLessThan(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "priority two", "--priority", "2"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "priority four", "--priority", "4"})

	exitCode := app.Run([]string{"list", "--priority", "3-"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "priority two") {
		t.Fatalf("expected priority 2 todo to be included, got %q", output)
	}

	if strings.Contains(output, "priority four") {
		t.Fatalf("expected priority 4 todo to be filtered out, got %q", output)
	}
}

func TestListFiltersPriorityNotEqual(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "priority three", "--priority", "3"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "priority four", "--priority", "4"})

	exitCode := app.Run([]string{"list", "--priority", "3!"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if strings.Contains(output, "priority three") {
		t.Fatalf("expected priority 3 todo to be filtered out, got %q", output)
	}

	if !strings.Contains(output, "priority four") {
		t.Fatalf("expected priority 4 todo to be included, got %q", output)
	}
}

func TestListSupportsExplicitFilters(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "first todo", "--status", "done", "--priority", "2"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "second todo", "--status", "done", "--priority", "4"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "third todo", "--status", "todo", "--priority", "4"})

	exitCode := app.Run([]string{"list", "--status", "done", "--priority", "3+"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if strings.Contains(output, "first todo") {
		t.Fatalf("expected first todo to be filtered out, got %q", output)
	}

	if !strings.Contains(output, "second todo") {
		t.Fatalf("expected matching todo to be included, got %q", output)
	}

	if strings.Contains(output, "third todo") {
		t.Fatalf("expected non-matching status todo to be filtered out, got %q", output)
	}
}

func TestListRejectsPositionalArgs(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"list", "done"})
	if exitCode != 1 {
		t.Fatalf("expected list to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "does not accept positional") {
		t.Fatalf("expected positional argument error, got %q", stderr.String())
	}
}

func TestDirectoryOptionIsParsedAfterCommand(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = filepath.Join(t.TempDir(), "default-missing")
	app.LookupEnv = func(string) (string, bool) {
		return "", false
	}
	app.Now = func() time.Time {
		return time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC)
	}

	otherDir := filepath.Join(t.TempDir(), "other")
	if err := os.Mkdir(otherDir, 0o750); err != nil {
		t.Fatalf("create other todo dir: %v", err)
	}

	exitCode := app.Run([]string{"add", "-d", otherDir, "--title", "Buy groceries"})
	if exitCode != 0 {
		t.Fatalf("expected add to succeed, got %d: %s", exitCode, stderr.String())
	}

	id := strings.TrimSpace(stdout.String())
	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"view", "-d", otherDir, id})
	if exitCode != 0 {
		t.Fatalf("expected view to succeed, got %d: %s", exitCode, stderr.String())
	}

	if !strings.Contains(stdout.String(), "title: Buy groceries") {
		t.Fatalf("expected view output from alternate directory, got %q", stdout.String())
	}
}

func TestDirectoryEnvironmentVariableIsUsed(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = filepath.Join(t.TempDir(), "default-missing")
	envDir := filepath.Join(t.TempDir(), "env-todo")
	app.LookupEnv = func(key string) (string, bool) {
		if key != "TODOLIST_DIRECTORY" {
			return "", false
		}

		return envDir, true
	}
	app.Now = func() time.Time {
		return time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC)
	}

	if err := os.Mkdir(envDir, 0o750); err != nil {
		t.Fatalf("create env todo dir: %v", err)
	}

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

	if !strings.Contains(stdout.String(), "title: Buy groceries") {
		t.Fatalf("expected view output from env directory, got %q", stdout.String())
	}
}

func TestInitCreatesTodoDirectoryAndConfig(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	todoDir := filepath.Join(t.TempDir(), "todo")
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = todoDir
	app.LookupEnv = func(string) (string, bool) {
		return "", false
	}

	exitCode := app.Run([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("expected init to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}

	info, err := os.Stat(todoDir)
	if err != nil {
		t.Fatalf("stat todo dir: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", todoDir)
	}

	configPath := filepath.Join(todoDir, ".todos")
	//nolint:gosec // Test reads a config file created in a temporary directory.
	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	if string(rawConfig) != "prefix=todo-\n" {
		t.Fatalf("expected default config contents, got %q", string(rawConfig))
	}

	output := stdout.String()
	if !strings.Contains(output, "initialized todo directory: "+todoDir) {
		t.Fatalf("expected initialization output, got %q", output)
	}

	if !strings.Contains(output, "created config file: "+configPath) {
		t.Fatalf("expected config creation output, got %q", output)
	}
}

func TestInitCreatesMissingConfigInExistingDirectory(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	todoDir := filepath.Join(t.TempDir(), "todo")
	if err := os.Mkdir(todoDir, 0o750); err != nil {
		t.Fatalf("create todo dir: %v", err)
	}

	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = todoDir
	app.LookupEnv = func(string) (string, bool) {
		return "", false
	}

	exitCode := app.Run([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("expected init to succeed, got %d: %s", exitCode, stderr.String())
	}

	configPath := filepath.Join(todoDir, ".todos")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("stat config: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "todo directory already exists: "+todoDir) {
		t.Fatalf("expected existing directory output, got %q", output)
	}

	if !strings.Contains(output, "created config file: "+configPath) {
		t.Fatalf("expected created config output, got %q", output)
	}
}

func TestInitIsIdempotentAndDoesNotOverwriteConfig(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	todoDir := filepath.Join(t.TempDir(), "todo")
	if err := os.Mkdir(todoDir, 0o750); err != nil {
		t.Fatalf("create todo dir: %v", err)
	}

	configPath := filepath.Join(todoDir, ".todos")
	if err := os.WriteFile(configPath, []byte("prefix=work-\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = todoDir
	app.LookupEnv = func(string) (string, bool) {
		return "", false
	}

	exitCode := app.Run([]string{"init"})
	if exitCode != 0 {
		t.Fatalf("expected init to succeed, got %d: %s", exitCode, stderr.String())
	}

	//nolint:gosec // Test reads a config file created in a temporary directory.
	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	if string(rawConfig) != "prefix=work-\n" {
		t.Fatalf("expected existing config to remain unchanged, got %q", string(rawConfig))
	}

	output := stdout.String()
	if !strings.Contains(output, "todo directory already exists: "+todoDir) {
		t.Fatalf("expected existing directory output, got %q", output)
	}

	if !strings.Contains(output, "config file already exists: "+configPath) {
		t.Fatalf("expected existing config output, got %q", output)
	}
}

func TestInitRespectsDirectoryOption(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = filepath.Join(t.TempDir(), "default-missing")
	app.LookupEnv = func(string) (string, bool) {
		return "", false
	}

	todoDir := filepath.Join(t.TempDir(), "work-todo")
	exitCode := app.Run([]string{"init", "-d", todoDir})
	if exitCode != 0 {
		t.Fatalf("expected init to succeed, got %d: %s", exitCode, stderr.String())
	}

	if _, err := os.Stat(todoDir); err != nil {
		t.Fatalf("stat initialized directory: %v", err)
	}

	if _, err := os.Stat(filepath.Join(todoDir, ".todos")); err != nil {
		t.Fatalf("stat initialized config: %v", err)
	}
}

func TestInitRejectsExtraPositionalArgs(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = filepath.Join(t.TempDir(), "todo")
	app.LookupEnv = func(string) (string, bool) {
		return "", false
	}

	exitCode := app.Run([]string{"init", "extra"})
	if exitCode != 1 {
		t.Fatalf("expected init to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "does not accept positional") {
		t.Fatalf("expected positional argument error, got %q", stderr.String())
	}
}

func TestInitRejectsNonDirectoryTodoPath(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	todoPath := filepath.Join(t.TempDir(), "todo")
	if err := os.WriteFile(todoPath, []byte("not a directory\n"), 0o600); err != nil {
		t.Fatalf("write todo path file: %v", err)
	}

	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = todoPath
	app.LookupEnv = func(string) (string, bool) {
		return "", false
	}

	exitCode := app.Run([]string{"init"})
	if exitCode != 1 {
		t.Fatalf("expected init to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `is not a directory`) {
		t.Fatalf("expected directory error, got %q", stderr.String())
	}
}

func TestInitRejectsNonRegularConfigPath(t *testing.T) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	todoDir := filepath.Join(t.TempDir(), "todo")
	if err := os.Mkdir(todoDir, 0o750); err != nil {
		t.Fatalf("create todo dir: %v", err)
	}

	configPath := filepath.Join(todoDir, ".todos")
	if err := os.Mkdir(configPath, 0o750); err != nil {
		t.Fatalf("create config dir: %v", err)
	}

	app := cli.NewApp(strings.NewReader(""), stdout, stderr, false)
	app.TodoDir = todoDir
	app.LookupEnv = func(string) (string, bool) {
		return "", false
	}

	exitCode := app.Run([]string{"init"})
	if exitCode != 1 {
		t.Fatalf("expected init to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `is not a regular file`) {
		t.Fatalf("expected config path error, got %q", stderr.String())
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

func TestUsagePrintsEmbeddedUsageText(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"usage"})
	if exitCode != 0 {
		t.Fatalf("expected usage to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, "## Using the todolist app") {
		t.Fatalf("expected usage heading, got %q", output)
	}

	if !strings.Contains(output, "todolist list --json") {
		t.Fatalf("expected usage content, got %q", output)
	}
}

func TestUsageHelpPrintsCommandUsage(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"usage", "--help"})
	if exitCode != 0 {
		t.Fatalf("expected usage help to succeed, got %d: %s", exitCode, stderr.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}

	if !strings.Contains(stdout.String(), "todolist usage") {
		t.Fatalf("expected usage command help, got %q", stdout.String())
	}
}

func TestUnknownCommandFails(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"nope"})
	if exitCode != 1 {
		t.Fatalf("expected unknown command to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `unknown command "nope"`) {
		t.Fatalf("expected unknown command error, got %q", stderr.String())
	}
}

func addTodoForTest(t *testing.T, app *cli.App, stdout, stderr *bytes.Buffer, args []string) string {
	t.Helper()

	stdout.Reset()
	stderr.Reset()
	app.StdinProvided = false
	app.Stdin = strings.NewReader("")

	exitCode := app.Run(args)
	if exitCode != 0 {
		t.Fatalf("expected %v to succeed, got %d: %s", args, exitCode, stderr.String())
	}

	id := strings.TrimSpace(stdout.String())
	stdout.Reset()
	stderr.Reset()

	return id
}

func newTestApp(t *testing.T, stdinProvided bool, stdin string) (*cli.App, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := cli.NewApp(strings.NewReader(stdin), stdout, stderr, stdinProvided)
	app.TodoDir = filepath.Join(t.TempDir(), "todo")
	app.UsageText = strings.Join([]string{
		"## Using the todolist app",
		"",
		"Use `todolist list --json` for structured output.",
		"",
	}, "\n")
	app.LookupEnv = func(string) (string, bool) {
		return "", false
	}
	app.Now = func() time.Time {
		return time.Date(2026, time.March, 18, 10, 0, 0, 0, time.UTC)
	}

	if err := os.Mkdir(app.TodoDir, 0o750); err != nil {
		t.Fatalf("create todo dir: %v", err)
	}

	return app, stdout, stderr
}
