package cli

import (
	"fmt"
	"strings"

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

	if _, err = app.Stdout.Write(raw); err != nil {
		return err
	}

	if len(value.Parents) == 0 {
		return nil
	}

	parentSection, err := renderParentSection(store, value.Parents)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(app.Stdout, "\n\n%s", parentSection)

	return err
}

func renderParentSection(store *todolist.Store, parentIDs []string) (string, error) {
	lines := []string{"Parents:"}
	for _, parentID := range parentIDs {
		parent, err := store.Get(parentID)
		if err != nil {
			return "", err
		}

		lines = append(lines, fmt.Sprintf("- %s %s", parent.ID, parent.Title))
	}

	return strings.Join(lines, "\n"), nil
}
