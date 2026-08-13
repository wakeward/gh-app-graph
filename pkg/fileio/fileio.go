// Package fileio centralizes safe local file writes for generated artifacts.
package fileio

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// PrivateFileMode is the default mode for generated files in this repo.
	PrivateFileMode = 0o600
	// PrivateDirMode is the default mode for directories we create.
	PrivateDirMode = 0o750
)

// Write writes data to path with PrivateFileMode.
func Write(path string, data []byte) error {
	if err := os.WriteFile(path, data, PrivateFileMode); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// MkdirAll creates path with PrivateDirMode.
func MkdirAll(path string) error {
	if err := os.MkdirAll(path, PrivateDirMode); err != nil {
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
