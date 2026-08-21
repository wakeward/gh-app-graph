package data

import "testing"

func TestLoadGHESOnlyAPIKeys(t *testing.T) {
	keys, err := LoadGHESOnlyAPIKeys()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"organization_pre_receive_hooks", "repository_pre_receive_hooks"}
	for _, key := range want {
		if _, ok := keys[key]; !ok {
			t.Errorf("expected GHES-only key %q in catalog", key)
		}
	}
	if len(keys) < len(want) {
		t.Fatalf("expected at least %d GHES-only keys, got %d", len(want), len(keys))
	}
}
