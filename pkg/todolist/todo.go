// Package todolist provides the core todo model, storage, and ID generation
// for the todolist CLI.
package todolist

import (
	"fmt"
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
	// CreatedAt is the todo creation timestamp in UTC.
	CreatedAt time.Time `json:"createdAt"`
	// LastModified is the timestamp of the most recent successful update in UTC.
	LastModified time.Time `json:"lastModified"`
	// Description is the raw Markdown description stored below the front matter.
	Description string `json:"description,omitempty"`
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

	value.CreatedAt = NormalizeTimestamp(value.CreatedAt)
	value.LastModified = NormalizeTimestamp(value.LastModified)

	return value
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
