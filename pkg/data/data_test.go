package data

import "testing"

func TestLoadToxicCombinations_EmbeddedCatalog(t *testing.T) {
	combos, err := LoadToxicCombinations()
	if err != nil {
		t.Fatal(err)
	}
	if len(combos) != 21 {
		t.Fatalf("expected 21 embedded toxic combinations, got %d", len(combos))
	}
}

func TestLoadResolvedPermissions_EmbeddedCatalog(t *testing.T) {
	perms, err := LoadResolvedPermissions()
	if err != nil {
		t.Fatal(err)
	}
	if len(perms) != 106 {
		t.Fatalf("expected 106 embedded permissions, got %d", len(perms))
	}
}
