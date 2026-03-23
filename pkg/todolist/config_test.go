package todolist_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kwo/todolist/pkg/todolist"
)

func TestLoadConfigDefaultsWhenFileIsMissing(t *testing.T) {
	t.Helper()

	config, err := todolist.LoadConfig(t.TempDir())
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if config.Prefix != todolist.DefaultIDPrefix {
		t.Fatalf("expected default prefix %q, got %q", todolist.DefaultIDPrefix, config.Prefix)
	}
}

func TestLoadConfigReadsConfiguredPrefix(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".todos"), []byte("prefix=work-\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := todolist.LoadConfig(dir)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if config.Prefix != "work-" {
		t.Fatalf("expected configured prefix %q, got %q", "work-", config.Prefix)
	}
}

func TestLoadConfigUsesDefaultPrefixWhenValueIsEmpty(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".todos"), []byte("prefix=\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := todolist.LoadConfig(dir)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if config.Prefix != todolist.DefaultIDPrefix {
		t.Fatalf("expected default prefix %q, got %q", todolist.DefaultIDPrefix, config.Prefix)
	}
}

func TestLoadConfigRejectsMalformedConfig(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".todos"), []byte("oops\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := todolist.LoadConfig(dir)
	if err == nil {
		t.Fatal("expected malformed config error")
	}
}
