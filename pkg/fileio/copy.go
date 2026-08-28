package fileio

import (
	"fmt"
	"os"
	"path/filepath"
)

// CopyFile copies src to dst with PrivateFileMode, creating parent dirs.
func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src) // #nosec G304 -- src is a generated artifact path from build
	if err != nil {
		return fmt.Errorf("read %s: %w", src, err)
	}
	if err := MkdirAll(filepath.Dir(dst)); err != nil {
		return err
	}
	return Write(dst, data)
}
