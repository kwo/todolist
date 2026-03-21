package todolist

import (
	"fmt"
	"hash/fnv"
	"strings"
)

const (
	alphabet  = "0123456789abcdefghjkmnpqrstvwxyz"
	prefix    = "todo-"
	visibleID = 4
)

// ExistsFunc reports whether a generated todo ID already exists.
type ExistsFunc func(string) bool

// GenerateID returns a unique todo ID for todo, retrying with incrementing nonce values on collision.
func GenerateID(todo Todo, exists ExistsFunc) string {
	for nonce := 0; ; nonce++ {
		id := prefix + suffix(buildSeed(todo, nonce))
		if !exists(id) {
			return id
		}
	}
}

func buildSeed(todo Todo, nonce int) string {
	createdAt := NormalizeTimestamp(todo.CreatedAt)

	return fmt.Sprintf("%s|%s|%d|%d", todo.Title, todo.Description, createdAt.Unix(), nonce)
}

func suffix(seed string) string {
	hash := fnv64a(seed)
	encoded := crockford(hash)
	if len(encoded) < visibleID {
		return strings.Repeat("0", visibleID-len(encoded)) + encoded
	}

	return encoded[len(encoded)-visibleID:]
}

func fnv64a(value string) uint64 {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(value))

	return hash.Sum64()
}

func crockford(value uint64) string {
	if value == 0 {
		return "0"
	}

	result := make([]byte, 0, 13)
	for value > 0 {
		result = append(result, alphabet[value%32])
		value /= 32
	}

	reverse(result)

	return string(result)
}

func reverse(value []byte) {
	for left, right := 0, len(value)-1; left < right; left, right = left+1, right-1 {
		value[left], value[right] = value[right], value[left]
	}
}
