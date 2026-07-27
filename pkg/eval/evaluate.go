// Package eval implements the toxic-combination matching engine: given a
// declared permission set, which named attack techniques does it enable.
//
// This is the package other tools should import directly. For example,
// gh-app-check's Control Plane auditor can call Evaluate with an
// installation's permissions instead of maintaining its own hardcoded
// toxic-permission predicates - see the root README's "Relationship to
// gh-app-check" section for the intended split of responsibilities.
package eval

import "github.com/wakeward/github-app-permissions-graph/pkg/model"

// Match is one toxic combination fully satisfied by a permission set.
type Match struct {
	Combination model.ToxicCombination
}

// NearMiss is a toxic combination missing exactly one permission grant to be
// fully satisfied - useful for "one permission away from X" warnings.
type NearMiss struct {
	Combination model.ToxicCombination
	Missing     []model.PermissionGrant
}

// Result is the full outcome of evaluating one permission set against every
// known toxic combination.
type Result struct {
	Matches      []Match
	NearMisses   []NearMiss
	HighestBlast model.BlastRadius
}

var blastRank = map[model.BlastRadius]int{
	"":                        0,
	model.BlastRadiusMedium:   1,
	model.BlastRadiusHigh:     2,
	model.BlastRadiusCritical: 3,
}

// Evaluate checks perms against every combination in combos and returns full
// matches and near-misses (combinations missing exactly one grant).
func Evaluate(perms model.AppPermissionSet, combos []model.ToxicCombination) Result {
	var result Result
	for _, combo := range combos {
		missing := missingGrants(perms, combo)
		switch len(missing) {
		case 0:
			result.Matches = append(result.Matches, Match{Combination: combo})
			if blastRank[combo.BlastRadius] > blastRank[result.HighestBlast] {
				result.HighestBlast = combo.BlastRadius
			}
		case 1:
			result.NearMisses = append(result.NearMisses, NearMiss{Combination: combo, Missing: missing})
		}
	}
	return result
}

func missingGrants(perms model.AppPermissionSet, combo model.ToxicCombination) []model.PermissionGrant {
	var missing []model.PermissionGrant
	for _, grant := range combo.Permissions {
		if !perms.HasGrant(grant.APIKey, grant.Access) {
			missing = append(missing, grant)
		}
	}
	return missing
}
