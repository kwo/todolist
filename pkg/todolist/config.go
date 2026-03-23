package todolist

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const configFileName = ".todos"

// Config contains per-directory todolist settings.
type Config struct {
	Prefix string
}

// LoadConfig reads the todo directory configuration, returning defaults when no config file exists.
func LoadConfig(dir string) (Config, error) {
	config := Config{Prefix: DefaultIDPrefix}
	path := filepath.Join(dir, configFileName)

	//nolint:gosec // The todo directory is an explicit local CLI input and .todos is a fixed filename within it.
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}

		return Config{}, fmt.Errorf("read todo config %q: %w", path, err)
	}

	parsed, err := parseConfig(string(raw))
	if err != nil {
		return Config{}, fmt.Errorf("parse todo config %q: %w", path, err)
	}

	return parsed, nil
}

func parseConfig(raw string) (Config, error) {
	config := Config{Prefix: DefaultIDPrefix}
	scanner := bufio.NewScanner(strings.NewReader(raw))

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return Config{}, fmt.Errorf("line %d: expected key=value", lineNumber)
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			return Config{}, fmt.Errorf("line %d: key cannot be empty", lineNumber)
		}

		switch key {
		case "prefix":
			if value == "" {
				config.Prefix = DefaultIDPrefix
				continue
			}

			config.Prefix = value
		default:
			return Config{}, fmt.Errorf("line %d: unknown key %q", lineNumber, key)
		}
	}

	if err := scanner.Err(); err != nil {
		return Config{}, fmt.Errorf("scan config: %w", err)
	}

	return config, nil
}
