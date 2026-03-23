package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kwo/todolist/pkg/cli"
)

//nolint:maintidx // End-to-end CLI flow assertions are intentionally kept together.
func TestAddListViewUpdateDeleteFlow(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, true, "Need milk, eggs, and bread.\n")

	exitCode := app.Run([]string{"add", "Buy groceries", "wip", "2"})
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

	expectedLine := id + "\tBuy groceries\n"
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

	exitCode = app.Run([]string{"update", id, "Buy groceries and snacks", "done", "1"})
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

func TestAddDefaultsStatusAndPriority(t *testing.T) {
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

	viewed := stdout.String()
	if !strings.Contains(viewed, "status: todo") {
		t.Fatalf("expected default status in view output, got %q", viewed)
	}

	if !strings.Contains(viewed, "priority: 5") {
		t.Fatalf("expected default priority in view output, got %q", viewed)
	}
}

func TestAddAllowsLiteralStatusLikeTitleWithExplicitAssignment(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "title=done"})
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

func TestAddAllowsLiteralPriorityLikeTitleWithExplicitAssignment(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "title=2"})
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

func TestAddRejectsInvalidExplicitStatus(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "status=active", "Buy groceries"})
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

func TestAddRejectsUnassignableExtraValue(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"add", "Buy groceries", "snacks"})
	if exitCode != 1 {
		t.Fatalf("expected add to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `cannot assign value "snacks"`) {
		t.Fatalf("expected assignment error, got %q", stderr.String())
	}
}

func TestUpdateAllowsLiteralStatusLikeTitleWithExplicitAssignment(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	id := addTodoForTest(t, app, stdout, stderr, []string{"add", "Buy groceries"})

	exitCode := app.Run([]string{"update", id, "title=done"})
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

func TestUpdateRejectsInvalidExplicitPriority(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	id := addTodoForTest(t, app, stdout, stderr, []string{"add", "Buy groceries"})

	exitCode := app.Run([]string{"update", id, "priority=6"})
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
	id := addTodoForTest(t, app, stdout, stderr, []string{"add", "Buy groceries"})

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

func TestListFiltersByStatus(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "first todo", "todo"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "second todo", "done"})

	exitCode := app.Run([]string{"list", "done"})
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

func TestListExcludesStatusWithBangPrefix(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "first todo", "todo"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "second todo", "done"})

	exitCode := app.Run([]string{"list", "!done"})
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

	addTodoForTest(t, app, stdout, stderr, []string{"add", "high priority", "1"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "low priority", "5"})

	exitCode := app.Run([]string{"list", "1"})
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

func TestListFiltersPriorityLessThan(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "priority two", "2"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "priority four", "4"})

	exitCode := app.Run([]string{"list", "-3"})
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

func TestListSupportsExplicitFilters(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "first todo", "done", "2"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "second todo", "done", "4"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "third todo", "todo", "4"})

	exitCode := app.Run([]string{"list", "status=done", "priority=+3"})
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

func TestListRejectsInvalidPriorityFilter(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"list", "."})
	if exitCode != 1 {
		t.Fatalf("expected list to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), `invalid priority filter ".": must use n, .n, +n, or -n`) {
		t.Fatalf("expected invalid priority filter error, got %q", stderr.String())
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

	exitCode := app.Run([]string{"add", "-d", otherDir, "Buy groceries"})
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
		t.Fatalf("expected view output from env directory, got %q", stdout.String())
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
