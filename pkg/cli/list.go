package cli

import (
	"fmt"
	"sort"

	"github.com/kwo/todolist/pkg/todolist"
)

type listCommand struct {
	StatusFilter   string
	ExcludeStatus  bool
	PriorityFilter string
	All            bool
}

func (c listCommand) Execute(app *App, options runOptions) error {
	priorityFilter, err := parsePriorityFilter(c.PriorityFilter)
	if err != nil {
		return err
	}

	todos, err := todolist.NewStore(options.TodoDir).List()
	if err != nil {
		return err
	}

	filtered := make([]todolist.Todo, 0, len(todos))
	for _, value := range todos {
		value = storeWithComputedFields(options.TodoDir, value)

		if c.StatusFilter != "" {
			matchesStatus := value.Status == c.StatusFilter
			if (!c.ExcludeStatus && !matchesStatus) || (c.ExcludeStatus && matchesStatus) {
				continue
			}
		}

		if priorityFilter != nil && !priorityFilter(value.Priority) {
			continue
		}

		if !c.All && !value.Ready {
			continue
		}

		filtered = append(filtered, value)
	}

	sortTodosForList(filtered)

	if options.JSON {
		return writeJSON(app.Stdout, filtered)
	}

	for _, value := range filtered {
		if _, err = fmt.Fprintf(app.Stdout, "%s\t%d\t%s\t%s\t%s\t%s\n", value.ID, value.Priority, value.Status, truncateListTitle(value.Title), formatListParents(value.Parents), formatListParents(value.Depends)); err != nil {
			return err
		}
	}

	return nil
}

func sortTodosForList(todos []todolist.Todo) {
	sort.Slice(todos, func(i, j int) bool {
		if todos[i].Priority != todos[j].Priority {
			return todos[i].Priority < todos[j].Priority
		}

		if todos[i].Title != todos[j].Title {
			return todos[i].Title < todos[j].Title
		}

		return todos[i].ID < todos[j].ID
	})
}

func truncateListTitle(title string) string {
	const maxTitleLength = 60
	const ellipsis = "..."

	runes := []rune(title)
	if len(runes) <= maxTitleLength {
		return title
	}

	return string(runes[:maxTitleLength-len([]rune(ellipsis))]) + ellipsis
}

func formatListParents(parents []string) string {
	if len(parents) == 0 {
		return ""
	}

	if len(parents) == 1 {
		return parents[0]
	}

	return parents[0] + ",..."
}
