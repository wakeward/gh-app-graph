package eval

import (
	"testing"

	"github.com/wakeward/gh-app-graph/pkg/model"
)

func stealthBackdoorCombo() model.ToxicCombination {
	return model.ToxicCombination{
		ID:          "stealth-backdoor",
		Technique:   "Stealth Backdoor",
		BlastRadius: model.BlastRadiusCritical,
		Permissions: []model.PermissionGrant{
			{APIKey: "administration", Access: model.AccessWrite},
			{APIKey: "contents", Access: model.AccessWrite},
		},
	}
}

func TestEvaluate_FullMatch(t *testing.T) {
	perms := model.AppPermissionSet{
		Permissions: map[string]string{
			"administration": "write",
			"contents":       "write",
		},
	}

	result := Evaluate(perms, []model.ToxicCombination{stealthBackdoorCombo()})

	if len(result.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(result.Matches))
	}
	if result.Matches[0].Combination.ID != "stealth-backdoor" {
		t.Errorf("expected stealth-backdoor match, got %q", result.Matches[0].Combination.ID)
	}
	if result.HighestBlast != model.BlastRadiusCritical {
		t.Errorf("expected HighestBlast=Critical, got %q", result.HighestBlast)
	}
	if len(result.NearMisses) != 0 {
		t.Errorf("expected 0 near misses, got %d", len(result.NearMisses))
	}
}

func TestEvaluate_NearMiss(t *testing.T) {
	perms := model.AppPermissionSet{
		Permissions: map[string]string{
			"administration": "write",
			// contents:write missing
		},
	}

	result := Evaluate(perms, []model.ToxicCombination{stealthBackdoorCombo()})

	if len(result.Matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(result.Matches))
	}
	if len(result.NearMisses) != 1 {
		t.Fatalf("expected 1 near miss, got %d", len(result.NearMisses))
	}
	missing := result.NearMisses[0].Missing
	if len(missing) != 1 || missing[0].APIKey != "contents" {
		t.Errorf("expected missing=[contents:write], got %+v", missing)
	}
}

func TestEvaluate_ReadOnlyDoesNotSatisfyWriteRequirement(t *testing.T) {
	perms := model.AppPermissionSet{
		Permissions: map[string]string{
			"administration": "write",
			"contents":       "read", // read does not satisfy a write requirement
		},
	}

	result := Evaluate(perms, []model.ToxicCombination{stealthBackdoorCombo()})

	if len(result.Matches) != 0 {
		t.Fatalf("expected 0 matches with contents:read only, got %d", len(result.Matches))
	}
	if len(result.NearMisses) != 1 {
		t.Fatalf("expected 1 near miss, got %d", len(result.NearMisses))
	}
}

func TestEvaluate_NoOverlapIsNeitherMatchNorNearMiss(t *testing.T) {
	perms := model.AppPermissionSet{
		Permissions: map[string]string{
			"issues": "write",
		},
	}

	result := Evaluate(perms, []model.ToxicCombination{stealthBackdoorCombo()})

	if len(result.Matches) != 0 || len(result.NearMisses) != 0 {
		t.Fatalf("expected no matches or near misses, got matches=%d nearMisses=%d", len(result.Matches), len(result.NearMisses))
	}
	if result.HighestBlast != "" {
		t.Errorf("expected empty HighestBlast, got %q", result.HighestBlast)
	}
}

func TestEvaluate_HighestBlastAcrossMultipleMatches(t *testing.T) {
	mediumCombo := model.ToxicCombination{
		ID:          "medium-combo",
		BlastRadius: model.BlastRadiusMedium,
		Permissions: []model.PermissionGrant{{APIKey: "issues", Access: model.AccessWrite}},
	}
	perms := model.AppPermissionSet{
		Permissions: map[string]string{
			"administration": "write",
			"contents":       "write",
			"issues":         "write",
		},
	}

	result := Evaluate(perms, []model.ToxicCombination{mediumCombo, stealthBackdoorCombo()})

	if len(result.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(result.Matches))
	}
	if result.HighestBlast != model.BlastRadiusCritical {
		t.Errorf("expected HighestBlast=Critical across matches, got %q", result.HighestBlast)
	}
}
