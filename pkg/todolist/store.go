package todolist

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

// Store reads and writes todo Markdown files in a single todo directory.
type Store struct {
	dir string
}

type frontMatter struct {
	ID           string   `yaml:"id"`
	Title        string   `yaml:"title"`
	Status       string   `yaml:"status"`
	Priority     int      `yaml:"priority"`
	Parents      []string `yaml:"parents,omitempty"`
	Depends      []string `yaml:"depends,omitempty"`
	CreatedAt    string   `yaml:"createdAt"`
	LastModified string   `yaml:"lastModified"`
}

// NewStore returns a Store rooted at dir.
func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// Exists reports whether a todo file already exists for id.
func (s *Store) Exists(id string) bool {
	_, err := os.Stat(s.pathFor(id))

	return err == nil
}

// Create writes a new todo file.
func (s *Store) Create(value Todo) error {
	if err := s.ensureDirectory(); err != nil {
		return err
	}

	value = NormalizeTodo(value)
	if err := s.validateTodo(value); err != nil {
		return err
	}

	return os.WriteFile(s.pathFor(value.ID), serialize(value), 0o600)
}

// Get loads a todo by ID.
func (s *Store) Get(id string) (Todo, error) {
	if err := s.ensureDirectory(); err != nil {
		return Todo{}, err
	}

	raw, err := os.ReadFile(s.pathFor(id))
	if err != nil {
		return Todo{}, fmt.Errorf("read todo %q: %w", id, err)
	}

	value, err := parse(raw)
	if err != nil {
		return Todo{}, err
	}

	if err := s.validateTodo(value); err != nil {
		return Todo{}, fmt.Errorf("read todo %q: %w", id, err)
	}

	return value, nil
}

// WithComputedFields returns a copy of value with computed output fields populated.
func (s *Store) WithComputedFields(value Todo) Todo {
	value.Ready = s.isReady(value.Depends)

	return value
}

// GetRaw loads the raw todo file bytes by ID.
func (s *Store) GetRaw(id string) ([]byte, error) {
	if err := s.ensureDirectory(); err != nil {
		return nil, err
	}

	raw, err := os.ReadFile(s.pathFor(id))
	if err != nil {
		return nil, fmt.Errorf("read todo %q: %w", id, err)
	}

	return raw, nil
}

// List loads all todo files in the store sorted by ID.
func (s *Store) List() ([]Todo, error) {
	if err := s.ensureDirectory(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("read todo directory %q: %w", s.dir, err)
	}

	todos := make([]Todo, 0, len(entries))
	for _, entry := range entries {
		if skipEntry(entry) {
			continue
		}

		item, itemErr := s.Get(strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name())))
		if itemErr != nil {
			return nil, itemErr
		}

		todos = append(todos, item)
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	return todos, nil
}

// Update overwrites an existing todo file.
func (s *Store) Update(value Todo) error {
	if err := s.ensureDirectory(); err != nil {
		return err
	}

	value = NormalizeTodo(value)
	if err := s.validateTodo(value); err != nil {
		return err
	}

	return os.WriteFile(s.pathFor(value.ID), serialize(value), 0o600)
}

// Delete removes a todo file by ID.
func (s *Store) Delete(id string) error {
	if err := s.ensureDirectory(); err != nil {
		return err
	}

	if err := os.Remove(s.pathFor(id)); err != nil {
		return fmt.Errorf("delete todo %q: %w", id, err)
	}

	return nil
}

func (s *Store) ensureDirectory() error {
	info, err := os.Stat(s.dir)
	if err != nil {
		return fmt.Errorf("todo directory %q: %w", s.dir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("todo directory %q is not a directory", s.dir)
	}

	return nil
}

func (s *Store) pathFor(id string) string {
	return filepath.Join(s.dir, id+".md")
}

func (s *Store) validateTodo(value Todo) error {
	if err := ValidateStatus(value.Status); err != nil {
		return err
	}

	if err := ValidatePriority(value.Priority); err != nil {
		return err
	}

	if err := ValidateParents(value.ID, value.Parents, s.Exists); err != nil {
		return err
	}

	if err := ValidateDepends(value.ID, value.Depends, s.Exists); err != nil {
		return err
	}

	return nil
}

func (s *Store) isReady(depends []string) bool {
	for _, dependencyID := range NormalizeDepends(depends) {
		dependency, err := s.getForReady(dependencyID)
		if err != nil {
			return false
		}

		if dependency.Status != "done" {
			return false
		}
	}

	return true
}

func (s *Store) getForReady(id string) (Todo, error) {
	if err := s.ensureDirectory(); err != nil {
		return Todo{}, err
	}

	raw, err := os.ReadFile(s.pathFor(id))
	if err != nil {
		return Todo{}, fmt.Errorf("read dependency todo %q: %w", id, err)
	}

	value, err := parse(raw)
	if err != nil {
		return Todo{}, err
	}

	if err := ValidateStatus(value.Status); err != nil {
		return Todo{}, err
	}

	if err := ValidatePriority(value.Priority); err != nil {
		return Todo{}, err
	}

	return value, nil
}

func parse(raw []byte) (Todo, error) {
	front, description, err := splitFrontMatter(string(raw))
	if err != nil {
		return Todo{}, err
	}

	value, err := parseFrontMatter(front)
	if err != nil {
		return Todo{}, err
	}

	createdAt, err := parseTime(value.CreatedAt)
	if err != nil {
		return Todo{}, err
	}

	lastModified, err := parseTime(value.LastModified)
	if err != nil {
		return Todo{}, err
	}

	return NormalizeTodo(Todo{
		ID:           value.ID,
		Title:        value.Title,
		Status:       value.Status,
		Priority:     value.Priority,
		Parents:      value.Parents,
		Depends:      value.Depends,
		CreatedAt:    createdAt,
		LastModified: lastModified,
		Description:  description,
	}), nil
}

func splitFrontMatter(raw string) (string, string, error) {
	if !strings.HasPrefix(raw, frontMatterBoundary) {
		return "", "", fmt.Errorf("todo file is missing YAML front matter")
	}

	rest := raw[len(frontMatterBoundary):]
	index := strings.Index(rest, frontMatterDivider)
	if index < 0 {
		return "", "", fmt.Errorf("todo file is missing a closing front matter boundary")
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

func serialize(value Todo) []byte {
	value = NormalizeTodo(value)

	metadata := frontMatter{
		ID:           value.ID,
		Title:        value.Title,
		Status:       value.Status,
		Priority:     value.Priority,
		Parents:      value.Parents,
		Depends:      value.Depends,
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
