package platform

import "github.com/wakeward/gh-app-graph/pkg/model"

// GHESOnlyKeys returns api_key values marked ghes_only in the permission catalog.
func GHESOnlyKeys(perms []model.Permission) map[string]struct{} {
	out := make(map[string]struct{})
	for _, p := range perms {
		if p.PlatformAvailability == model.PlatformGHESOnly {
			out[p.APIKey] = struct{}{}
		}
	}
	return out
}

// ComboAvailability returns the platform requirement for a toxic combination.
func ComboAvailability(combo model.ToxicCombination, ghesKeys map[string]struct{}) model.PlatformAvailability {
	if combo.PlatformAvailability != "" {
		return combo.PlatformAvailability
	}
	for _, grant := range combo.Permissions {
		if _, ok := ghesKeys[grant.APIKey]; ok {
			return model.PlatformGHESOnly
		}
	}
	return model.PlatformAll
}

// FilterToxicCombinations returns combos that apply on the scan target.
// When includeGHES is false, GHES-only rules are excluded.
func FilterToxicCombinations(combos []model.ToxicCombination, ghesKeys map[string]struct{}, includeGHES bool) []model.ToxicCombination {
	if includeGHES {
		return combos
	}
	filtered := make([]model.ToxicCombination, 0, len(combos))
	for _, combo := range combos {
		if ComboAvailability(combo, ghesKeys) == model.PlatformGHESOnly {
			continue
		}
		filtered = append(filtered, combo)
	}
	return filtered
}

// GrantedGHESOnly lists GHES-only scopes present in a permission map.
func GrantedGHESOnly(perms map[string]string, ghesKeys map[string]struct{}) []string {
	if len(ghesKeys) == 0 || len(perms) == 0 {
		return nil
	}
	var out []string
	for key, access := range perms {
		if _, ok := ghesKeys[key]; ok {
			out = append(out, key+":"+access)
		}
	}
	return out
}
