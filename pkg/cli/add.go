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
	Parents  []string
	Depends  []string
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
		Parents:      c.Parents,
		Depends:      c.Depends,
		CreatedAt:    now,
		LastModified: now,
		Description:  description,
	}
	value.ID = todolist.GenerateIDWithPrefix(value, config.Prefix, store.Exists)

	if err := todolist.ValidateParents(value.ID, value.Parents, store.Exists); err != nil {
		return err
	}

	if err := todolist.ValidateDepends(value.ID, value.Depends, store.Exists); err != nil {
		return err
	}

	if err := store.Create(value); err != nil {
		return err
	}

	value = todolist.NormalizeTodo(value)

	if options.JSON {
		value = store.WithComputedFields(value)

		return writeJSON(app.Stdout, value)
	}

	_, err = fmt.Fprintf(app.Stdout, "%s\n", value.ID)

	return err
}
