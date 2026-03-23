package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kwo/todolist/pkg/todolist"
)

type listCommand struct {
	StatusFilter   string
	ExcludeStatus  bool
	PriorityFilter string
}

func parseListCommand(values []string) (listCommand, error) {
	command := listCommand{}
	statusAssigned := false
	priorityAssigned := false

	for _, raw := range values {
		key, value, explicit := parseAssignment(raw)
		if explicit {
			switch key {
			case "status":
				if statusAssigned {
					return listCommand{}, fmt.Errorf("status filter was provided more than once")
				}

				statusFilter, excludeStatus, err := parseExplicitStatusFilter(value)
				if err != nil {
					return listCommand{}, err
				}

				command.StatusFilter = statusFilter
				command.ExcludeStatus = excludeStatus
				statusAssigned = true
			case "priority":
				if priorityAssigned {
					return listCommand{}, fmt.Errorf("priority filter was provided more than once")
				}

				if _, err := parsePriorityFilter(value); err != nil {
					return listCommand{}, err
				}

				command.PriorityFilter = strings.TrimSpace(value)
				priorityAssigned = true
			default:
				return listCommand{}, fmt.Errorf("unknown list field %q", key)
			}

			continue
		}

		if !statusAssigned {
			statusFilter, excludeStatus, ok, err := recognizeStatusFilter(raw)
			if err != nil {
				return listCommand{}, err
			}

			if ok {
				command.StatusFilter = statusFilter
				command.ExcludeStatus = excludeStatus
				statusAssigned = true

				continue
			}
		}

		if !priorityAssigned {
			priorityFilter, ok, err := recognizePriorityFilter(raw)
			if err != nil {
				return listCommand{}, err
			}

			if ok {
				command.PriorityFilter = priorityFilter
				priorityAssigned = true

				continue
			}
		}

		return listCommand{}, fmt.Errorf("cannot assign value %q", raw)
	}

	if !statusAssigned {
		command.StatusFilter = "done"
		command.ExcludeStatus = true
	}

	return command, nil
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

func parseStatusFilter(raw string) (string, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, nil
	}

	exclude := strings.HasSuffix(value, "!")
	if exclude {
		value = strings.TrimSpace(strings.TrimSuffix(value, "!"))
	}

	if err := todolist.ValidateStatus(value); err != nil {
		return "", false, err
	}

	return value, exclude, nil
}

func parseExplicitStatusFilter(raw string) (string, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, fmt.Errorf("invalid status %q: must be one of todo, wip, done", value)
	}

	return parseStatusFilter(value)
}

func parsePriorityFilter(raw string) (func(int) bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return nil, nil
	}

	operator := "="
	switch {
	case hasPriorityFilterPrefix(value):
		operator = value[:1]
		value = strings.TrimSpace(value[1:])
	case hasPriorityFilterSuffix(value):
		operator = value[len(value)-1:]
		value = strings.TrimSpace(value[:len(value)-1])
	}

	priority, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid priority filter %q: must use n, n!, n+, or n-", raw)
	}

	if priority < 0 || priority > todolist.DefaultPriority {
		return nil, fmt.Errorf("invalid priority filter %q: priority must be between 0 and %d", raw, todolist.DefaultPriority)
	}

	switch operator {
	case "=":
		return func(candidate int) bool { return candidate == priority }, nil
	case ".", "!":
		return func(candidate int) bool { return candidate != priority }, nil
	case "+":
		return func(candidate int) bool { return candidate > priority }, nil
	case "-":
		return func(candidate int) bool { return candidate < priority }, nil
	default:
		return nil, fmt.Errorf("invalid priority filter %q", raw)
	}
}

func recognizeStatusFilter(raw string) (string, bool, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, false, nil
	}

	if strings.HasSuffix(value, "!") {
		status, exclude, err := parseStatusFilter(value)
		if err == nil {
			return status, exclude, true, nil
		}
	}

	if strings.HasPrefix(value, "!") {
		status, exclude, err := parseStatusFilter(value)

		return status, exclude, true, err
	}

	status, ok := recognizeStatusValue(value)
	if !ok {
		return "", false, false, nil
	}

	return status, false, true, nil
}

func recognizePriorityFilter(raw string) (string, bool, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, nil
	}

	if !isPriorityFilterShape(value) {
		return "", false, nil
	}

	if _, err := parsePriorityFilter(value); err != nil {
		return "", true, err
	}

	return value, true, nil
}

func isPriorityFilterShape(value string) bool {
	if isDigits(value) {
		return true
	}

	if len(value) < 2 {
		return value == "." || value == "+" || value == "-" || value == "!"
	}

	if hasPriorityFilterPrefix(value) {
		return isDigits(strings.TrimSpace(value[1:])) || strings.TrimSpace(value[1:]) == ""
	}

	if hasPriorityFilterSuffix(value) {
		return isDigits(strings.TrimSpace(value[:len(value)-1])) || strings.TrimSpace(value[:len(value)-1]) == ""
	}

	return false
}

func hasPriorityFilterPrefix(value string) bool {
	switch value[0] {
	case '.', '+', '-':
		return true
	default:
		return false
	}
}

func hasPriorityFilterSuffix(value string) bool {
	switch value[len(value)-1] {
	case '!', '+', '-':
		return true
	default:
		return false
	}
}
