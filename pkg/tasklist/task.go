// Package tasklist provides the core task model, storage, and ID generation
// for the tasklist CLI.
package tasklist

import (
	"fmt"
	"time"
)

const (
	// DefaultStatus is the status used when a task omits a status value.
	DefaultStatus = "todo"
	// DefaultPriority is the priority used when a task omits a priority value.
	DefaultPriority = 5
)

// Task is a single task stored as a Markdown file with YAML front matter.
type Task struct {
	// ID is the stable task identifier and filename prefix.
	ID string
	// Title is the human-readable task title.
	Title string
	// Status is the workflow state of the task.
	Status string
	// Priority is the task priority where 1 is highest and 5 is lowest.
	Priority int
	// CreatedAt is the task creation timestamp in UTC.
	CreatedAt time.Time
	// LastModified is the timestamp of the most recent successful update in UTC.
	LastModified time.Time
	// Description is the raw Markdown description stored below the front matter.
	Description string
}

// NormalizeTimestamp converts a timestamp to UTC and truncates it to whole-second precision.
func NormalizeTimestamp(value time.Time) time.Time {
	return value.UTC().Truncate(time.Second)
}

// NormalizeTask normalizes stored task fields to their canonical representation.
func NormalizeTask(value Task) Task {
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

// ValidateStatus reports whether status is one of the supported task status values.
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
