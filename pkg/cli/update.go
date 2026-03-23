package cli

import (
	"fmt"
	"strings"

	"github.com/kwo/todolist/pkg/todolist"
)

type updateCommand struct {
	Todo             string
	Title            string
	Status           string
	Priority         int
	TitleProvided    bool
	StatusProvided   bool
	PriorityProvided bool
}

func parseUpdateCommand(values []string) (updateCommand, error) {
	if len(values) == 0 {
		return updateCommand{}, fmt.Errorf("update requires a todo id")
	}

	command := updateCommand{Todo: values[0]}
	titleAssigned := false
	statusAssigned := false
	priorityAssigned := false

	for _, raw := range values[1:] {
		key, value, explicit := parseAssignment(raw)
		if explicit {
			switch key {
			case "title":
				if titleAssigned {
					return updateCommand{}, fmt.Errorf("title was provided more than once")
				}

				command.Title = value
				command.TitleProvided = true
				titleAssigned = true
			case "status":
				if statusAssigned {
					return updateCommand{}, fmt.Errorf("status was provided more than once")
				}

				if err := todolist.ValidateStatus(strings.TrimSpace(value)); err != nil {
					return updateCommand{}, err
				}

				command.Status = strings.TrimSpace(value)
				command.StatusProvided = true
				statusAssigned = true
			case "priority":
				if priorityAssigned {
					return updateCommand{}, fmt.Errorf("priority was provided more than once")
				}

				priority, err := parseExplicitPriority(value)
				if err != nil {
					return updateCommand{}, err
				}

				command.Priority = priority
				command.PriorityProvided = true
				priorityAssigned = true
			default:
				return updateCommand{}, fmt.Errorf("unknown update field %q", key)
			}

			continue
		}

		if !statusAssigned {
			status, ok := recognizeStatusValue(raw)
			if ok {
				command.Status = status
				command.StatusProvided = true
				statusAssigned = true

				continue
			}
		}

		if !priorityAssigned {
			priority, ok := recognizePriorityValue(raw)
			if ok {
				command.Priority = priority
				command.PriorityProvided = true
				priorityAssigned = true

				continue
			}
		}

		if !titleAssigned {
			command.Title = raw
			command.TitleProvided = true
			titleAssigned = true

			continue
		}

		return updateCommand{}, fmt.Errorf("cannot assign value %q", raw)
	}

	return command, nil
}

func (c updateCommand) Execute(app *App, options runOptions) error {
	description, descriptionProvided, err := readOptionalDescription(app.Stdin, app.StdinProvided)
	if err != nil {
		return err
	}

	store := todolist.NewStore(options.TodoDir)
	value, err := store.Get(c.Todo)
	if err != nil {
		return err
	}

	if c.TitleProvided {
		value.Title = strings.TrimSpace(c.Title)
		if value.Title == "" {
			return fmt.Errorf("title cannot be empty")
		}
	}

	if c.StatusProvided {
		value.Status = c.Status
	}

	if c.PriorityProvided {
		value.Priority = c.Priority
	}

	if descriptionProvided {
		value.Description = description
	}

	if !c.TitleProvided && !c.StatusProvided && !c.PriorityProvided && !descriptionProvided {
		return fmt.Errorf("update requires a title, status, priority, or stdin description input")
	}

	value.LastModified = todolist.NormalizeTimestamp(app.Now())

	if err := store.Update(value); err != nil {
		return err
	}

	if options.JSON {
		return writeJSON(app.Stdout, value)
	}

	return nil
}
