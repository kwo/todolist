package cli_test

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestListIncludesIDPriorityStatusAndTitleColumns(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	id := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Buy groceries", "--status", "wip", "--priority", "2"})

	exitCode := app.Run([]string{"list"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	expectedLine := formatListLine(id, 2, "wip", "Buy groceries", "", "")
	if stdout.String() != expectedLine {
		t.Fatalf("expected list output %q, got %q", expectedLine, stdout.String())
	}
}

func TestListSortsByPriorityThenTitle(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	lowPriority := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Alpha", "--priority", "4"})
	priorityOneZulu := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Zulu", "--priority", "1"})
	priorityOneAlpha := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Alpha", "--priority", "1"})

	exitCode := app.Run([]string{"list"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	expected := strings.Join([]string{
		strings.TrimSuffix(formatListLine(priorityOneAlpha, 1, "todo", "Alpha", "", ""), "\n"),
		strings.TrimSuffix(formatListLine(priorityOneZulu, 1, "todo", "Zulu", "", ""), "\n"),
		strings.TrimSuffix(formatListLine(lowPriority, 4, "todo", "Alpha", "", ""), "\n"),
	}, "\n") + "\n"
	if stdout.String() != expected {
		t.Fatalf("expected sorted list output %q, got %q", expected, stdout.String())
	}
}

func TestListJSONSortsByPriorityThenTitle(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Bravo", "--priority", "2"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Alpha", "--priority", "2"})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Zulu", "--priority", "1"})

	exitCode := app.Run([]string{"list", "--json"})
	if exitCode != 0 {
		t.Fatalf("expected json list to succeed, got %d: %s", exitCode, stderr.String())
	}

	var listed []jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &listed); err != nil {
		t.Fatalf("unmarshal list json: %v; output=%q", err, stdout.String())
	}

	if len(listed) != 3 {
		t.Fatalf("expected 3 todos, got %+v", listed)
	}

	titles := []string{listed[0].Title, listed[1].Title, listed[2].Title}
	if strings.Join(titles, ",") != "Zulu,Alpha,Bravo" {
		t.Fatalf("expected json list order by priority then title, got %+v", titles)
	}

	priorities := []int{listed[0].Priority, listed[1].Priority, listed[2].Priority}
	if priorities[0] != 1 || priorities[1] != 2 || priorities[2] != 2 {
		t.Fatalf("expected json priorities [1 2 2], got %+v", priorities)
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

	truncatedTitle := truncateListTitleForTest(longTitle)
	if len(truncatedTitle) != 60 {
		t.Fatalf("expected truncated title length 60, got %d in %q", len(truncatedTitle), truncatedTitle)
	}

	if !strings.HasSuffix(truncatedTitle, "...") {
		t.Fatalf("expected truncated title to end with ellipsis, got %q", truncatedTitle)
	}

	if strings.Contains(truncatedTitle, "vendors") {
		t.Fatalf("expected truncated title not to include full title, got %q", truncatedTitle)
	}

	expectedLine := formatListLine(id, 3, "todo", truncatedTitle, "", "")
	if stdout.String() != expectedLine {
		t.Fatalf("expected list output %q, got %q", expectedLine, stdout.String())
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

func TestListTextOutputHasFixedWidth(t *testing.T) {
	t.Helper()

	line := formatListLine("todo-7k9m", 5, "todo", strings.Repeat("x", 60), "todo-1a2b,...", "todo-3c4d,...")
	line = strings.TrimSuffix(line, "\n")

	if len(line) != 110 {
		t.Fatalf("expected fixed-width line length 110, got %d in %q", len(line), line)
	}
}

func TestListDefaultShowsOnlyReadyTodos(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	readyDependency := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Ready dep", "--status", "done"})
	blockedDependency := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Blocked dep", "--status", "wip"})
	readyID := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Ready todo", "--depends", readyDependency})
	blockedID := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Blocked todo", "--depends", readyDependency, "--depends", blockedDependency})

	exitCode := app.Run([]string{"list"})
	if exitCode != 0 {
		t.Fatalf("expected list to succeed, got %d: %s", exitCode, stderr.String())
	}

	if !strings.Contains(stdout.String(), strings.TrimSuffix(formatListLine(readyID, 5, "todo", "Ready todo", "", readyDependency), "\n")) {
		t.Fatalf("expected ready todo in default list output, got %q", stdout.String())
	}

	if strings.Contains(stdout.String(), blockedID) {
		t.Fatalf("expected blocked todo to be excluded by default, got %q", stdout.String())
	}
}

func TestListAllIncludesReadyAndBlockedTodos(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	readyDependency := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Ready dep", "--status", "done"})
	blockedDependency := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Blocked dep", "--status", "wip"})
	readyID := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Ready todo", "--depends", readyDependency})
	blockedID := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Blocked todo", "--depends", readyDependency, "--depends", blockedDependency})

	exitCode := app.Run([]string{"list", "--all"})
	if exitCode != 0 {
		t.Fatalf("expected list --all to succeed, got %d: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, strings.TrimSuffix(formatListLine(readyID, 5, "todo", "Ready todo", "", readyDependency), "\n")) {
		t.Fatalf("expected ready todo in --all output, got %q", output)
	}

	if !strings.Contains(output, strings.TrimSuffix(formatListLine(blockedID, 5, "todo", "Blocked todo", "", readyDependency+",..."), "\n")) {
		t.Fatalf("expected blocked todo in --all output, got %q", output)
	}

	stdout.Reset()
	stderr.Reset()

	exitCode = app.Run([]string{"list", "--json", "--all"})
	if exitCode != 0 {
		t.Fatalf("expected json list --all to succeed, got %d: %s", exitCode, stderr.String())
	}

	var listed []jsonTodo
	if err := json.Unmarshal(stdout.Bytes(), &listed); err != nil {
		t.Fatalf("unmarshal list json: %v; output=%q", err, stdout.String())
	}

	var sawReady bool
	var sawBlocked bool
	for _, todo := range listed {
		switch todo.ID {
		case readyID:
			sawReady = true
		case blockedID:
			sawBlocked = true
			if len(todo.Depends) != 2 || todo.Depends[0] != readyDependency || todo.Depends[1] != blockedDependency {
				t.Fatalf("expected depends in json output, got %+v", todo)
			}
			if todo.Ready {
				t.Fatalf("expected blocked todo to have ready=false, got %+v", todo)
			}
		}
	}

	if !sawReady || !sawBlocked {
		t.Fatalf("expected both ready and blocked todos in json output, got %+v", listed)
	}
}

func TestListRejectsRemovedReadyFlag(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")

	exitCode := app.Run([]string{"list", "--ready", "maybe"})
	if exitCode != 1 {
		t.Fatalf("expected removed ready flag to fail, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}

	if !strings.Contains(stderr.String(), "unknown flag: --ready") {
		t.Fatalf("expected unknown ready flag error, got %q", stderr.String())
	}
}

func TestListAllComposesWithStatusAndPriority(t *testing.T) {
	t.Helper()

	app, stdout, stderr := newTestApp(t, false, "")
	readyDependency := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Ready dep", "--status", "done"})
	blockedDependency := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Blocked dep", "--status", "wip"})
	blockedID := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Blocked wip", "--status", "wip", "--priority", "2", "--depends", blockedDependency})
	matchingID := addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Ready wip", "--status", "wip", "--priority", "2", "--depends", readyDependency})
	addTodoForTest(t, app, stdout, stderr, []string{"add", "--title", "Wrong priority", "--status", "wip", "--priority", "4", "--depends", readyDependency})

	exitCode := app.Run([]string{"list", "--status", "wip", "--priority", "2", "--all"})
	if exitCode != 0 {
		t.Fatalf("expected composed list filters with --all to succeed, got %d: %s", exitCode, stderr.String())
	}

	output := stdout.String()
	if !strings.Contains(output, strings.TrimSuffix(formatListLine(matchingID, 2, "wip", "Ready wip", "", readyDependency), "\n")) {
		t.Fatalf("expected matching ready todo, got %q", output)
	}

	if !strings.Contains(output, strings.TrimSuffix(formatListLine(blockedID, 2, "wip", "Blocked wip", "", blockedDependency), "\n")) {
		t.Fatalf("expected matching blocked todo, got %q", output)
	}

	if strings.Contains(output, "Wrong priority") {
		t.Fatalf("expected only priority 2 matches, got %q", output)
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
