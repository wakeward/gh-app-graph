package graph

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wakeward/gh-app-graph/pkg/fileio"
	"github.com/wakeward/gh-app-graph/pkg/model"
)

func TestMergeDependencies_Precedence(t *testing.T) {
	overlap := []model.DependencyEdge{{
		Permission: "workflows",
		Access:     model.AccessWrite,
		Requires: []model.RequiredPermission{{
			Permission: "contents",
			Access:     model.AccessWrite,
			Type:       model.DependencyTypeCategoryOverlap,
		}},
	}}
	scraped := []model.DependencyEdge{{
		Permission: "workflows",
		Access:     model.AccessWrite,
		Requires: []model.RequiredPermission{{
			Permission: "contents",
			Access:     model.AccessRead,
			Type:       model.DependencyTypeProseScraped,
		}},
	}}
	manual := []model.DependencyEdge{{
		Permission: "workflows",
		Access:     model.AccessWrite,
		Requires: []model.RequiredPermission{{
			Permission: "contents",
			Access:     model.AccessWrite,
			Type:       model.DependencyTypeManual,
			Condition:  "human confirmed",
		}},
	}}

	merged := MergeDependencies(overlap, scraped, manual)
	if len(merged) != 1 {
		t.Fatalf("expected 1 edge, got %d", len(merged))
	}
	if merged[0].Requires[0].Type != model.DependencyTypeManual {
		t.Fatalf("expected manual to win, got %q", merged[0].Requires[0].Type)
	}
}

func TestLoadPermissionsDir_RejectsDuplicateAPIKey(t *testing.T) {
	dir := t.TempDir()
	entry := "- name: A\n  api_key: dup\n  category: account\n  overview: x\n  access_levels: []\n  needs_investigation: false\n  doc_status: documented\n  impact_plane: control\n"
	for _, name := range []string{"a.yaml", "b.yaml"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(entry), fileio.PrivateFileMode); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := LoadPermissionsDir(dir); err == nil {
		t.Fatal("expected duplicate api_key error")
	}
}
