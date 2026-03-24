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
