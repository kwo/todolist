// Package todolist provides the core todo model, storage, and ID generation
// for the todolist CLI.
package todolist

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

const (
	// DefaultStatus is the status used when a todo omits a status value.
	DefaultStatus = "todo"
	// DefaultPriority is the priority used when a todo omits a priority value.
	DefaultPriority = 5
)

// Todo is a single todo stored as a Markdown file with YAML front matter.
type Todo struct {
	// ID is the stable todo identifier and filename prefix.
	ID string `json:"id"`
	// Title is the human-readable todo title.
	Title string `json:"title"`
	// Status is the workflow state of the todo.
	Status string `json:"status"`
	// Priority is the todo priority where 1 is highest and 5 is lowest.
	Priority int `json:"priority"`
	// Parents contains zero or more parent todo IDs used for grouping.
	Parents []string `json:"parents,omitempty"`
	// Depends contains zero or more todo IDs that must be done before this todo is ready.
	Depends []string `json:"depends,omitempty"`
	// CreatedAt is the todo creation timestamp in UTC.
	CreatedAt time.Time `json:"createdAt"`
	// LastModified is the timestamp of the most recent successful update in UTC.
	LastModified time.Time `json:"lastModified"`
	// Description is the raw Markdown description stored below the front matter.
	Description string `json:"description,omitempty"`
	// Ready reports whether all dependencies are done. It is computed, not stored.
	Ready bool `json:"ready"`
}

// NormalizeTimestamp converts a timestamp to UTC and truncates it to whole-second precision.
func NormalizeTimestamp(value time.Time) time.Time {
	return value.UTC().Truncate(time.Second)
}

// NormalizeTodo normalizes stored todo fields to their canonical representation.
func NormalizeTodo(value Todo) Todo {
	if value.Status == "" {
		value.Status = DefaultStatus
	}

	if value.Priority == 0 {
		value.Priority = DefaultPriority
	}

	value.Parents = NormalizeParents(value.Parents)
	value.Depends = NormalizeDepends(value.Depends)
	value.CreatedAt = NormalizeTimestamp(value.CreatedAt)
	value.LastModified = NormalizeTimestamp(value.LastModified)

	return value
}

// NormalizeParents trims whitespace, drops empty values, and keeps stable order.
func NormalizeParents(parents []string) []string {
	return normalizeReferences(parents, false)
}

// NormalizeDepends trims whitespace, drops empty values, deduplicates, and keeps stable order.
func NormalizeDepends(depends []string) []string {
	return normalizeReferences(depends, true)
}

func normalizeReferences(values []string, dedupe bool) []string {
	if len(values) == 0 {
		return nil
	}

	normalized := make([]string, 0, len(values))
	seen := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		if dedupe {
			if _, ok := seen[value]; ok {
				continue
			}

			seen[value] = struct{}{}
		}

		normalized = append(normalized, value)
	}

	if len(normalized) == 0 {
		return nil
	}

	return normalized
}

// ValidateStatus reports whether status is one of the supported todo status values.
func ValidateStatus(status string) error {
	switch status {
	case "todo", "wip", "done":
		return nil
	default:
		return fmt.Errorf("invalid status %q: must be one of todo, wip, done", status)
	}
}

// ValidatePriority reports whether priority is one of the supported priority values.
func ValidatePriority(priority int) error {
	if priority < 1 || priority > 5 {
		return fmt.Errorf("invalid priority %d: must be between 1 and 5", priority)
	}

	return nil
}

// ValidateParents validates parent todo IDs for duplicates, self-references, and existence.
func ValidateParents(id string, parents []string, exists func(string) bool) error {
	normalized := NormalizeParents(parents)
	seen := make(map[string]struct{}, len(normalized))
	for _, parent := range normalized {
		if parent == id {
			return fmt.Errorf("invalid parent %q: a todo cannot be its own parent", parent)
		}

		if _, ok := seen[parent]; ok {
			return fmt.Errorf("duplicate parent %q", parent)
		}

		seen[parent] = struct{}{}

		if exists != nil && !exists(parent) {
			return fmt.Errorf("parent todo %q does not exist", parent)
		}
	}

	if !slices.Equal(normalized, parents) {
		// Empty or whitespace-only parent IDs are treated as invalid when supplied explicitly.
		for _, parent := range parents {
			if strings.TrimSpace(parent) == "" {
				return fmt.Errorf("invalid parent %q", parent)
			}
		}
	}

	return nil
}

// ValidateDepends validates dependency todo IDs for self-references and existence.
func ValidateDepends(id string, depends []string, exists func(string) bool) error {
	normalized := NormalizeDepends(depends)
	for _, dependency := range normalized {
		if dependency == id {
			return fmt.Errorf("invalid dependency %q: a todo cannot depend on itself", dependency)
		}

		if exists != nil && !exists(dependency) {
			return fmt.Errorf("dependency todo %q does not exist", dependency)
		}
	}

	if !slices.Equal(normalized, depends) {
		for _, dependency := range depends {
			if strings.TrimSpace(dependency) == "" {
				return fmt.Errorf("invalid dependency %q", dependency)
			}
		}
	}

	return nil
}
