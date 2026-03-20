package tasklist

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	frontMatterBoundary = "---\n"
	frontMatterDivider  = "\n---\n"
)

// Store reads and writes task Markdown files in a single task directory.
type Store struct {
	dir string
}

type frontMatter struct {
	ID           string `yaml:"id"`
	Title        string `yaml:"title"`
	Status       string `yaml:"status"`
	Priority     int    `yaml:"priority"`
	CreatedAt    string `yaml:"createdAt"`
	LastModified string `yaml:"lastModified"`
}

// NewStore returns a Store rooted at dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Exists reports whether a task file already exists for id.
func (s *Store) Exists(id string) bool {
	_, err := os.Stat(s.pathFor(id))

	return err == nil
}

// Create writes a new task file.
func (s *Store) Create(value Task) error {
	if err := s.ensureDirectory(); err != nil {
		return err
	}

	return os.WriteFile(s.pathFor(value.ID), serialize(NormalizeTask(value)), 0o600)
}

// Get loads a task by ID.
func (s *Store) Get(id string) (Task, error) {
	if err := s.ensureDirectory(); err != nil {
		return Task{}, err
	}

	raw, err := os.ReadFile(s.pathFor(id))
	if err != nil {
		return Task{}, fmt.Errorf("read task %q: %w", id, err)
	}

	return parse(raw)
}

// GetRaw loads the raw task file bytes by ID.
func (s *Store) GetRaw(id string) ([]byte, error) {
	if err := s.ensureDirectory(); err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(s.pathFor(id))
	if err != nil {
		return nil, fmt.Errorf("read task %q: %w", id, err)
	}

	return raw, nil
}

// List loads all task files in the store sorted by ID.
func (s *Store) List() ([]Task, error) {
	if err := s.ensureDirectory(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("read task directory %q: %w", s.dir, err)
	}

	tasks := make([]Task, 0, len(entries))
	for _, entry := range entries {
		if skipEntry(entry) {
			continue
		}

		item, itemErr := s.Get(strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())))
		if itemErr != nil {
			return nil, itemErr
		}

		tasks = append(tasks, item)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks, nil
}

// Update overwrites an existing task file.
func (s *Store) Update(value Task) error {
	if err := s.ensureDirectory(); err != nil {
		return err
	}

	return os.WriteFile(s.pathFor(value.ID), serialize(NormalizeTask(value)), 0o600)
}

// Delete removes a task file by ID.
func (s *Store) Delete(id string) error {
	if err := s.ensureDirectory(); err != nil {
		return err
	}

	if err := os.Remove(s.pathFor(id)); err != nil {
		return fmt.Errorf("delete task %q: %w", id, err)
	}

	return nil
}

func (s *Store) ensureDirectory() error {
	info, err := os.Stat(s.dir)
	if err != nil {
		return fmt.Errorf("task directory %q: %w", s.dir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("task directory %q is not a directory", s.dir)
	}

	return nil
}

func (s *Store) pathFor(id string) string {
	return filepath.Join(s.dir, id+".md")
}

func parse(raw []byte) (Task, error) {
	front, description, err := splitFrontMatter(string(raw))
	if err != nil {
		return Task{}, err
	}

	value, err := parseFrontMatter(front)
	if err != nil {
		return Task{}, err
	}

	createdAt, err := parseTime(value.CreatedAt)
	if err != nil {
		return Task{}, err
	}

	lastModified, err := parseTime(value.LastModified)
	if err != nil {
		return Task{}, err
	}

	return NormalizeTask(Task{
		ID:           value.ID,
		Title:        value.Title,
		Status:       value.Status,
		Priority:     value.Priority,
		CreatedAt:    createdAt,
		LastModified: lastModified,
		Description:  description,
	}), nil
}

func splitFrontMatter(raw string) (string, string, error) {
	if !strings.HasPrefix(raw, frontMatterBoundary) {
		return "", "", fmt.Errorf("task file is missing YAML front matter")
	}

	rest := raw[len(frontMatterBoundary):]
	index := strings.Index(rest, frontMatterDivider)
	if index < 0 {
		return "", "", fmt.Errorf("task file is missing a closing front matter boundary")
	}

	front := rest[:index]
	description := strings.TrimPrefix(rest[index+len(frontMatterDivider):], "\n")

	return front, description, nil
}

func parseFrontMatter(raw string) (frontMatter, error) {
	value := frontMatter{}
	if err := yaml.Unmarshal([]byte(raw), &value); err != nil {
		return frontMatter{}, fmt.Errorf("parse front matter: %w", err)
	}

	return value, nil
}

func parseTime(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse time %q: %w", value, err)
	}

	return parsed.UTC(), nil
}

func serialize(value Task) []byte {
	value = NormalizeTask(value)

	metadata := frontMatter{
		ID:           value.ID,
		Title:        value.Title,
		Status:       value.Status,
		Priority:     value.Priority,
		CreatedAt:    formatTime(value.CreatedAt),
		LastModified: formatTime(value.LastModified),
	}

	raw, err := yaml.Marshal(metadata)
	if err != nil {
		panic(err)
	}

	buffer := bytes.Buffer{}
	buffer.WriteString(frontMatterBoundary)
	buffer.Write(raw)
	buffer.WriteString(frontMatterBoundary)
	buffer.WriteString("\n")
	buffer.WriteString(value.Description)

	return buffer.Bytes()
}

func formatTime(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func skipEntry(entry os.DirEntry) bool {
	return entry.IsDir() || filepath.Ext(entry.Name()) != ".md"
}
