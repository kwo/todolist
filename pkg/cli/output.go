package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/kwo/todolist/pkg/todolist"
)

func readDescription(reader io.Reader, provided bool) (string, error) {
	value, _, err := readOptionalDescription(reader, provided)

	return value, err
}

func readOptionalDescription(reader io.Reader, provided bool) (string, bool, error) {
	if !provided {
		return "", false, nil
	}

	raw, err := io.ReadAll(reader)
	if err != nil {
		return "", false, fmt.Errorf("read stdin: %w", err)
	}

	return string(raw), true, nil
}

func writeJSON(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)

	return encoder.Encode(value)
}

func storeWithComputedFields(todoDir string, value todolist.Todo) todolist.Todo {
	return todolist.NewStore(todoDir).WithComputedFields(value)
}
