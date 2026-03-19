package main

import (
	"os"

	"github.com/kwo/tasklist/pkg/cli"
)

func main() {
	app := cli.NewApp(os.Stdin, os.Stdout, os.Stderr, stdinProvided())
	os.Exit(app.Run(os.Args[1:]))
}

func stdinProvided() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice == 0
}
