package todolist

import (
	"errors"
	"fmt"
	"strings"
)

const (
	alphabet        = "0123456789abcdefghjkmnpqrstvwxyz"
	DefaultIDPrefix = "todo-"
	visibleID       = 4
	idRadix         = len(alphabet)
	maxVisibleIDs   = idRadix * idRadix * idRadix * idRadix
	maxCounter      = maxVisibleIDs - 1
)

var (
	// ErrIDSpaceExhausted reports that no four-character IDs are available.
	ErrIDSpaceExhausted = errors.New("todo id space exhausted")
	// ErrInvalidLastID reports an invalid last_id configuration value.
	ErrInvalidLastID = errors.New("invalid last todo id")
)

// ExistsFunc reports whether a generated todo ID already exists.
type ExistsFunc func(string) bool

// GenerateID returns the next available todo ID using the default prefix.
func GenerateID(lastID string, exists ExistsFunc) (string, error) {
	return GenerateIDWithPrefix(lastID, DefaultIDPrefix, exists)
}

// GenerateIDWithPrefix returns the next available todo ID for prefix.
func GenerateIDWithPrefix(lastID, prefix string, exists ExistsFunc) (string, error) {
	counter, err := nextCounter(lastID, prefix)
	if err != nil {
		return "", err
	}

	for ; counter < maxVisibleIDs; counter++ {
		id := prefix + encodeCounter(counter)
		if !exists(id) {
			return id, nil
		}
	}

	return "", ErrIDSpaceExhausted
}

func nextCounter(lastID, prefix string) (int, error) {
	trimmedLastID := strings.TrimSpace(lastID)
	if trimmedLastID == "" {
		return 0, nil
	}

	if !strings.HasPrefix(trimmedLastID, prefix) {
		return 0, fmt.Errorf("%w: %q does not start with prefix %q", ErrInvalidLastID, trimmedLastID, prefix)
	}

	suffix := strings.TrimPrefix(trimmedLastID, prefix)
	if len(suffix) != visibleID {
		return 0, fmt.Errorf("%w: %q has suffix length %d, expected %d", ErrInvalidLastID, trimmedLastID, len(suffix), visibleID)
	}

	counter, err := decodeCounter(suffix)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalidLastID, err)
	}

	if counter == maxCounter {
		return maxVisibleIDs, nil
	}

	return counter + 1, nil
}

func decodeCounter(suffix string) (int, error) {
	value := 0
	for i := 0; i < len(suffix); i++ {
		digit := strings.IndexByte(alphabet, suffix[i])
		if digit < 0 {
			return 0, fmt.Errorf("suffix %q contains unsupported character %q", suffix, string(suffix[i]))
		}

		value = value*idRadix + digit
	}

	return value, nil
}

func encodeCounter(value int) string {
	result := make([]byte, visibleID)
	for i := visibleID - 1; i >= 0; i-- {
		result[i] = alphabet[value%idRadix]
		value /= idRadix
	}

	return string(result)
}
