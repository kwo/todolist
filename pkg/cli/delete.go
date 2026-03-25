package cli

import "github.com/kwo/todolist/pkg/todolist"

type deleteCommand viewCommand

type deleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

func (c deleteCommand) Execute(app *App, options runOptions) error {
	store := todolist.NewStore(options.TodoDir)
	todos, err := store.List()
	if err != nil {
		return err
	}

	for _, value := range todos {
		if !slicesContains(value.Parents, c.Todo) {
			continue
		}

		value.Parents = removeParent(value.Parents, c.Todo)
		value.LastModified = todolist.NormalizeTimestamp(app.Now())
		if err := store.Update(value); err != nil {
			return err
		}
	}

	if err := store.Delete(c.Todo); err != nil {
		return err
	}

	if options.JSON {
		return writeJSON(app.Stdout, deleteResult{ID: c.Todo, Deleted: true})
	}

	return nil
}
