package cli

import (
	"fmt"
	"strings"

	"github.com/kwo/todolist/pkg/todolist"
)

type updateCommand struct {
	Todo             string
	Title            string
	Status           string
	Priority         int
	Parents          []string
	Depends          []string
	TitleProvided    bool
	StatusProvided   bool
	PriorityProvided bool
	ParentsProvided  bool
	DependsProvided  bool
}

func (c updateCommand) Execute(app *App, options runOptions) error {
	description, descriptionProvided, err := readOptionalDescription(app.Stdin, app.StdinProvided)
	if err != nil {
		return err
	}

	store := todolist.NewStore(options.TodoDir)
	value, err := store.Get(c.Todo)
	if err != nil {
		return err
	}

	if c.TitleProvided {
		value.Title = strings.TrimSpace(c.Title)
		if value.Title == "" {
			return fmt.Errorf("title cannot be empty")
		}
	}

	if c.StatusProvided {
		value.Status = c.Status
	}

	if c.PriorityProvided {
		value.Priority = c.Priority
	}

	originalParents := append([]string(nil), value.Parents...)

	if c.ParentsProvided {
		updatedParents, parentErr := applyParentOperations(value.ID, value.Parents, c.Parents, store.Exists)
		if parentErr != nil {
			return parentErr
		}

		value.Parents = updatedParents
	}

	if c.DependsProvided {
		updatedDepends, dependencyErr := applyDependencyOperations(value.ID, value.Depends, c.Depends, store.Exists)
		if dependencyErr != nil {
			return dependencyErr
		}

		value.Depends = updatedDepends
	}

	if descriptionProvided {
		value.Description = description
	}

	if !c.TitleProvided && !c.StatusProvided && !c.PriorityProvided && !c.ParentsProvided && !c.DependsProvided && !descriptionProvided {
		return fmt.Errorf("update requires a title, status, priority, parent, dependency, or stdin description input")
	}

	value.LastModified = todolist.NormalizeTimestamp(app.Now())

	if err := store.Update(value); err != nil {
		return err
	}

	if c.ParentsProvided {
		removedParents, addedParents := diffParents(originalParents, value.Parents)
		if err := syncParentDependencyLinks(store, value.ID, removedParents, addedParents, value.LastModified); err != nil {
			return err
		}
	}

	if options.JSON {
		value = store.WithComputedFields(value)

		return writeJSON(app.Stdout, value)
	}

	return nil
}

func applyParentOperations(todoID string, current, operations []string, exists func(string) bool) ([]string, error) {
	parents := append([]string(nil), todolist.NormalizeParents(current)...)
	seenOps := map[string]string{}

	for _, operation := range operations {
		raw := strings.TrimSpace(operation)
		if raw == "" {
			return nil, fmt.Errorf("invalid parent %q", operation)
		}

		action := "add"
		parentID := raw
		if strings.HasSuffix(raw, "!") {
			action = "remove"
			parentID = strings.TrimSpace(strings.TrimSuffix(raw, "!"))
		}

		if parentID == "" {
			return nil, fmt.Errorf("invalid parent %q", operation)
		}

		if prior, ok := seenOps[parentID]; ok {
			if prior != action {
				return nil, fmt.Errorf("conflicting parent operations for %q", parentID)
			}

			return nil, fmt.Errorf("duplicate parent operation for %q", parentID)
		}

		seenOps[parentID] = action

		switch action {
		case "add":
			if err := todolist.ValidateParents(todoID, []string{parentID}, exists); err != nil {
				return nil, err
			}

			if slicesContains(parents, parentID) {
				return nil, fmt.Errorf("duplicate parent %q", parentID)
			}

			parents = append(parents, parentID)
		case "remove":
			if parentID == todoID {
				return nil, fmt.Errorf("invalid parent %q: a todo cannot be its own parent", parentID)
			}

			if !slicesContains(parents, parentID) {
				return nil, fmt.Errorf("parent %q is not currently assigned", parentID)
			}

			parents = removeParent(parents, parentID)
		}
	}

	return parents, nil
}

func applyDependencyOperations(todoID string, current, operations []string, exists func(string) bool) ([]string, error) {
	depends := append([]string(nil), todolist.NormalizeDepends(current)...)
	seenOps := map[string]string{}

	for _, operation := range operations {
		raw := strings.TrimSpace(operation)
		if raw == "" {
			return nil, fmt.Errorf("invalid dependency %q", operation)
		}

		action := "add"
		dependencyID := raw
		if strings.HasSuffix(raw, "!") {
			action = "remove"
			dependencyID = strings.TrimSpace(strings.TrimSuffix(raw, "!"))
		}

		if dependencyID == "" {
			return nil, fmt.Errorf("invalid dependency %q", operation)
		}

		if prior, ok := seenOps[dependencyID]; ok {
			if prior != action {
				return nil, fmt.Errorf("conflicting dependency operations for %q", dependencyID)
			}

			if action == "remove" {
				return nil, fmt.Errorf("duplicate dependency operation for %q", dependencyID)
			}

			continue
		}

		seenOps[dependencyID] = action

		switch action {
		case "add":
			if err := todolist.ValidateDepends(todoID, []string{dependencyID}, exists); err != nil {
				return nil, err
			}

			if slicesContains(depends, dependencyID) {
				continue
			}

			depends = append(depends, dependencyID)
		case "remove":
			if dependencyID == todoID {
				return nil, fmt.Errorf("invalid dependency %q: a todo cannot depend on itself", dependencyID)
			}

			if !slicesContains(depends, dependencyID) {
				return nil, fmt.Errorf("dependency %q is not currently assigned", dependencyID)
			}

			depends = removeParent(depends, dependencyID)
		}
	}

	return depends, nil
}

func slicesContains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}

	return false
}

func removeParent(values []string, needle string) []string {
	filtered := values[:0]
	for _, value := range values {
		if value != needle {
			filtered = append(filtered, value)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	return filtered
}
