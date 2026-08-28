package fileio

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateRelativePath_RejectsTraversal(t *testing.T) {
	for _, path := range []string{"../outside", "../../etc/passwd", "/tmp/abs"} {
		if _, err := validateRelativePath(path); err == nil {
			t.Fatalf("expected error for %q", path)
		}
	}
}

func TestWrite_AllowsRepoRelativePath(t *testing.T) {
	dir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
	t.Cleanup(func() { _ = os.Chdir(wd) })

	rel := filepath.Join("generated", "out.txt")
	if err := MkdirAll(filepath.Dir(rel)); err != nil {
		t.Fatal(err)
	}
	if err := Write(rel, []byte("ok")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(rel); err != nil {
		t.Fatalf("expected file: %v", err)
	}
}
