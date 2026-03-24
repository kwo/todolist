package cli

import (
	"fmt"

	"github.com/kwo/todolist/pkg/todolist"
)

type listCommand struct {
	StatusFilter   string
	ExcludeStatus  bool
	PriorityFilter string
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
		if c.StatusFilter != "" {
			matchesStatus := value.Status == c.StatusFilter
			if (!c.ExcludeStatus && !matchesStatus) || (c.ExcludeStatus && matchesStatus) {
				continue
			}
		}

		if priorityFilter != nil && !priorityFilter(value.Priority) {
			continue
		}

		filtered = append(filtered, value)
	}

	if options.JSON {
		return writeJSON(app.Stdout, filtered)
	}

	for _, value := range filtered {
		if _, err = fmt.Fprintf(app.Stdout, "%s\t%d\t%s\t%s\n", value.ID, value.Priority, value.Status, truncateListTitle(value.Title)); err != nil {
			return err
		}
	}

	return nil
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
