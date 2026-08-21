package platform

import (
	"testing"

	"github.com/wakeward/gh-app-graph/pkg/model"
)

func TestFilterToxicCombinations_ExcludesGHESOnCloud(t *testing.T) {
	ghesKeys := map[string]struct{}{"repository_pre_receive_hooks": {}}
	combos := []model.ToxicCombination{
		{ID: "cloud-combo", PlatformAvailability: model.PlatformAll},
		{ID: "ghes-combo", Permissions: []model.PermissionGrant{
			{APIKey: "repository_pre_receive_hooks", Access: model.AccessWrite},
			{APIKey: "contents", Access: model.AccessWrite},
		}},
	}

	filtered := FilterToxicCombinations(combos, ghesKeys, false)
	if len(filtered) != 1 || filtered[0].ID != "cloud-combo" {
		t.Fatalf("expected only cloud combo, got %+v", filtered)
	}
}

func TestGrantedGHESOnly(t *testing.T) {
	ghesKeys := map[string]struct{}{
		"organization_pre_receive_hooks": {},
	}
	grants := map[string]string{
		"contents":                       "write",
		"organization_pre_receive_hooks": "read",
	}
	got := GrantedGHESOnly(grants, ghesKeys)
	if len(got) != 1 || got[0] != "organization_pre_receive_hooks:read" {
		t.Fatalf("got %v", got)
	}
}
