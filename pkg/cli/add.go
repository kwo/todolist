package cli

import (
	"fmt"
	"strings"

	"github.com/kwo/todolist/pkg/todolist"
)

type addCommand struct {
	Title    string
	Status   string
	Priority int
}

func parseAddCommand(values []string) (addCommand, error) {
	command := addCommand{
		Status:   todolist.DefaultStatus,
		Priority: todolist.DefaultPriority,
	}

	titleAssigned := false
	statusAssigned := false
	priorityAssigned := false

	for _, raw := range values {
		key, value, explicit := parseAssignment(raw)
		if explicit {
			switch key {
			case "title":
				if titleAssigned {
					return addCommand{}, fmt.Errorf("title was provided more than once")
				}

				command.Title = value
				titleAssigned = true
			case "status":
				if statusAssigned {
					return addCommand{}, fmt.Errorf("status was provided more than once")
				}

				if err := todolist.ValidateStatus(strings.TrimSpace(value)); err != nil {
					return addCommand{}, err
				}

				command.Status = strings.TrimSpace(value)
				statusAssigned = true
			case "priority":
				if priorityAssigned {
					return addCommand{}, fmt.Errorf("priority was provided more than once")
				}

				priority, err := parseExplicitPriority(value)
				if err != nil {
					return addCommand{}, err
				}

				command.Priority = priority
				priorityAssigned = true
			default:
				return addCommand{}, fmt.Errorf("unknown add field %q", key)
			}

			continue
		}

		if !statusAssigned {
			status, ok := recognizeStatusValue(raw)
			if ok {
				command.Status = status
				statusAssigned = true

				continue
			}
		}

		if !priorityAssigned {
			priority, ok := recognizePriorityValue(raw)
			if ok {
				command.Priority = priority
				priorityAssigned = true

				continue
			}
		}

		if !titleAssigned {
			command.Title = raw
			titleAssigned = true

			continue
		}

		return addCommand{}, fmt.Errorf("cannot assign value %q", raw)
	}

	if !titleAssigned {
		return addCommand{}, fmt.Errorf("add requires a title")
	}

	return command, nil
}

func (c addCommand) Execute(app *App, options runOptions) error {
	title := strings.TrimSpace(c.Title)
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}

	description, err := readDescription(app.Stdin, app.StdinProvided)
	if err != nil {
		return err
	}

	now := todolist.NormalizeTimestamp(app.Now())
	store := todolist.NewStore(options.TodoDir)
	config, err := todolist.LoadConfig(options.TodoDir)
	if err != nil {
		return err
	}

	value := todolist.Todo{
		Title:        title,
		Status:       c.Status,
		Priority:     c.Priority,
		CreatedAt:    now,
		LastModified: now,
		Description:  description,
	}
	value.ID = todolist.GenerateIDWithPrefix(value, config.Prefix, store.Exists)

	if err := store.Create(value); err != nil {
		return err
	}

	if options.JSON {
		return writeJSON(app.Stdout, value)
	}

	_, err = fmt.Fprintf(app.Stdout, "%s\n", value.ID)

	return err
}
