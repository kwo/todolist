package cli

import "fmt"

const defaultUsageText = "Usage documentation is not embedded in this build.\n"

type usageCommand struct{}

func (c usageCommand) Execute(app *App, options runOptions) error {
	text := app.UsageText
	if text == "" {
		text = defaultUsageText
	}

	_, err := fmt.Fprint(app.Stdout, text)

	return err
}
