// Package fileio centralizes safe local file writes for generated artifacts.
package fileio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// PrivateFileMode is the default mode for generated files in this repo.
	PrivateFileMode = 0o600
	// PrivateDirMode is the default mode for directories we create.
	PrivateDirMode = 0o750
)

// validateRelativePath rejects absolute paths and traversal outside the working directory.
func validateRelativePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("absolute paths are not allowed: %s", path)
	}
	clean := filepath.Clean(path)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes working directory: %s", path)
	}
	return clean, nil
}

func anchoredPath(path string) (string, error) {
	rel, err := validateRelativePath(path)
	if err != nil {
		return "", err
	}
	root, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	abs := filepath.Clean(filepath.Join(root, rel))
	root = filepath.Clean(root)
	if abs != root && !strings.HasPrefix(abs, root+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes working directory: %s", path)
	}
	return abs, nil
}

// Write writes data to path with PrivateFileMode.
func Write(path string, data []byte) error {
	target, err := anchoredPath(path)
	if err != nil {
		return err
	}
	if err := os.WriteFile(target, data, PrivateFileMode); err != nil { // #nosec G703 -- target from anchoredPath under cwd
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// MkdirAll creates path with PrivateDirMode.
func MkdirAll(path string) error {
	target, err := anchoredPath(path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(target, PrivateDirMode); err != nil { // #nosec G703 -- target from anchoredPath under cwd
		return fmt.Errorf("mkdir %s: %w", path, err)
	}
	return nil
}

// WriteYAML writes a YAML header plus encoded value to path.
func WriteYAML(path string, header string, v any) error {
	var buf strings.Builder
	if header != "" {
		buf.WriteString(header)
	}
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode yaml: %w", err)
	}
	if err := enc.Close(); err != nil {
		return fmt.Errorf("close yaml encoder: %w", err)
	}
	return Write(path, []byte(buf.String()))
}
