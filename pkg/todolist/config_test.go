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

	if config.LastID != "" {
		t.Fatalf("expected empty default last id, got %q", config.LastID)
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

	if config.LastID != "" {
		t.Fatalf("expected empty last id when not configured, got %q", config.LastID)
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

func TestLoadConfigReadsConfiguredLastID(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ".todos"), []byte("prefix=work-\nlast_id=work-000f\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := todolist.LoadConfig(dir)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if config.Prefix != "work-" {
		t.Fatalf("expected configured prefix %q, got %q", "work-", config.Prefix)
	}

	if config.LastID != "work-000f" {
		t.Fatalf("expected configured last id %q, got %q", "work-000f", config.LastID)
	}
}

func TestSaveConfigWritesPrefixAndLastID(t *testing.T) {
	t.Helper()

	dir := t.TempDir()
	config := todolist.Config{Prefix: "work-", LastID: "work-000f"}

	if err := todolist.SaveConfig(dir, config); err != nil {
		t.Fatalf("save config: %v", err)
	}

	//nolint:gosec // Test reads a config file from a temporary directory.
	raw, err := os.ReadFile(filepath.Join(dir, ".todos"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	if string(raw) != "prefix=work-\nlast_id=work-000f\n" {
		t.Fatalf("expected serialized config, got %q", string(raw))
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
