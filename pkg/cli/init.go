package cli

import (
	"fmt"

	"github.com/kwo/todolist/pkg/todolist"
)

type initCommand struct{}

func (c initCommand) Execute(app *App, options runOptions) error {
	result, err := todolist.InitDirectory(options.TodoDir)
	if err != nil {
		return err
	}

	if options.JSON {
		return writeJSON(app.Stdout, result)
	}

	if result.DirectoryCreated {
		if _, err := fmt.Fprintf(app.Stdout, "initialized todo directory: %s\n", result.Directory); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(app.Stdout, "todo directory already exists: %s\n", result.Directory); err != nil {
			return err
		}
	}

	if result.ConfigCreated {
		if _, err := fmt.Fprintf(app.Stdout, "created config file: %s\n", result.ConfigPath); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(app.Stdout, "config file already exists: %s\n", result.ConfigPath); err != nil {
			return err
		}
	}

	return nil
}
