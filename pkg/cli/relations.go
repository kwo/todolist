package cli

import (
	"time"

	"github.com/kwo/todolist/pkg/todolist"
)

func syncParentDependencyLinks(store *todolist.Store, childID string, removedParents, addedParents []string, now time.Time) error {
	timestamp := todolist.NormalizeTimestamp(now)

	for _, parentID := range todolist.NormalizeParents(removedParents) {
		parent, err := store.Get(parentID)
		if err != nil {
			return err
		}

		if !slicesContains(parent.Depends, childID) {
			continue
		}

		parent.Depends = removeParent(parent.Depends, childID)
		parent.LastModified = timestamp
		if err := store.Update(parent); err != nil {
			return err
		}
	}

	for _, parentID := range todolist.NormalizeParents(addedParents) {
		parent, err := store.Get(parentID)
		if err != nil {
			return err
		}

		if slicesContains(parent.Depends, childID) {
			continue
		}

		parent.Depends = append(parent.Depends, childID)
		parent.LastModified = timestamp
		if err := store.Update(parent); err != nil {
			return err
		}
	}

	return nil
}

func diffParents(current, updated []string) (removed, added []string) {
	current = todolist.NormalizeParents(current)
	updated = todolist.NormalizeParents(updated)

	for _, parentID := range current {
		if !slicesContains(updated, parentID) {
			removed = append(removed, parentID)
		}
	}

	for _, parentID := range updated {
		if !slicesContains(current, parentID) {
			added = append(added, parentID)
		}
	}

	return removed, added
}
