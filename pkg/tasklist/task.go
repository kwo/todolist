// Package tasklist provides the core task model, storage, and ID generation
// for the tasklist CLI.
package tasklist

import "time"

// Task is a single task stored as a Markdown file with YAML front matter.
type Task struct {
	// ID is the stable task identifier and filename prefix.
	ID string
	// Title is the human-readable task title.
	Title string
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

// NormalizeTask normalizes all timestamps on a task to the stored precision.
func NormalizeTask(value Task) Task {
	value.CreatedAt = NormalizeTimestamp(value.CreatedAt)
	value.LastModified = NormalizeTimestamp(value.LastModified)

	return value
}
