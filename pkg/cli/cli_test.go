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
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Status       string   `json:"status"`
	Priority     int      `json:"priority"`
	Parents      []string `json:"parents"`
	Depends      []string `json:"depends"`
	CreatedAt    string   `json:"createdAt"`
	LastModified string   `json:"lastModified"`
	Description  string   `json:"description"`
	Ready        bool     `json:"ready"`
}

type jsonDeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
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

	expectedLine := id + "\t2\twip\tready\tBuy groceries\t\t\n"
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

	if viewed.ID != added.ID || viewed.Title != added.Title || viewed.Status != added.Status || viewed.Priority != added.Priority || viewed.CreatedAt != added.CreatedAt || viewed.LastModified != added.LastModified || viewed.Description != added.Description || strings.Join(viewed.Parents, ",") != strings.Join(added.Parents, ",") {
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
