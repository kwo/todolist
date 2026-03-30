package cli

import (
	"github.com/kwo/todolist/pkg/todolist"
)

type viewCommand struct {
	Todo string
}

func (c viewCommand) Execute(app *App, options runOptions) error {
	store := todolist.NewStore(options.TodoDir)
	value, err := store.Get(c.Todo)
	if err != nil {
		return err
	}

	if options.JSON {
		value = store.WithComputedFields(value)

		return writeJSON(app.Stdout, value)
	}

	raw, err := store.GetRaw(c.Todo)
	if err != nil {
		return err
	}

	_, err = app.Stdout.Write(raw)

	return err
}
