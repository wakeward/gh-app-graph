package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/wakeward/gh-app-graph/pkg/model"
)

var writeupComboRE = regexp.MustCompile(`\[\s*([^:\]]+?)\s*:\s*(read|write)\s*\]\s*\+\s*\[\s*([^:\]]+?)\s*:\s*(read|write)\s*\]`)

// supplementFromWriteup parses the "Summary Matrix: Toxic Combinations" section
// of the threat writeup and returns combos not already present in existing.
func supplementFromWriteup(path string, existing []model.ToxicCombination, known map[string]struct{}) ([]model.ToxicCombination, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path is a caller-supplied writeup for one-time migration
	if err != nil {
		return nil, err
	}
	text := strings.ReplaceAll(string(data), `\`, "")
	section := text
	if i := strings.Index(strings.ToLower(text), "summary matrix"); i >= 0 {
		section = text[i:]
	}

	seen := make(map[string]struct{}, len(existing))
	for _, c := range existing {
		seen[grantSetKey(c.Permissions)] = struct{}{}
	}

	var added []model.ToxicCombination
	for _, m := range writeupComboRE.FindAllStringSubmatch(section, -1) {
		leftKey, err := resolveAPIKey(m[1], known)
		if err != nil {
			return nil, fmt.Errorf("writeup left grant %q: %w", m[1], err)
		}
		rightKey, err := resolveAPIKey(m[3], known)
		if err != nil {
			return nil, fmt.Errorf("writeup right grant %q: %w", m[3], err)
		}
		leftAccess, err := parseAccessLevel(m[2])
		if err != nil {
			return nil, err
		}
		rightAccess, err := parseAccessLevel(m[4])
		if err != nil {
			return nil, err
		}
		grants := []model.PermissionGrant{
			{APIKey: leftKey, Access: leftAccess},
			{APIKey: rightKey, Access: rightAccess},
		}
		key := grantSetKey(grants)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		technique := techniqueForWriteupPair(leftKey, rightKey)
		added = append(added, model.ToxicCombination{
			ID:              slugify(technique),
			Technique:       technique,
			Permissions:     grants,
			BlastRadius:     blastForWriteupPair(leftKey, rightKey),
			RiskDescription: writeupRiskDescription(leftKey, rightKey),
			ExploitPath:     writeupExploitPath(section, leftKey, rightKey),
			Source:          "migrate-toxic-combinations writeup supplement",
		})
	}
	return added, nil
}

func techniqueForWriteupPair(a, b string) string {
	pair := grantSetKey([]model.PermissionGrant{
		{APIKey: a, Access: model.AccessWrite},
		{APIKey: b, Access: model.AccessWrite},
	})
	switch pair {
	case "administration:write+contents:write":
		return "Stealth Backdoor"
	case "contents:write+workflows:write":
		return "Arbitrary Code Execution"
	case "checks:write+pull_requests:write":
		return "Status Check Forgery Merge"
	case "members:write+organization_administration:write":
		return "Organization Takeover"
	default:
		return fmt.Sprintf("%s + %s", a, b)
	}
}

func blastForWriteupPair(a, b string) model.BlastRadius {
	switch a + "+" + b {
	case "checks+pull_requests", "pull_requests+checks":
		return model.BlastRadiusHigh
	default:
		return model.BlastRadiusCritical
	}
}

func writeupRiskDescription(a, b string) string {
	switch {
	case (a == "checks" && b == "pull_requests") || (a == "pull_requests" && b == "checks"):
		return "Checks: write can forge required status results; Pull requests: write can open and merge changes without human review when combined."
	case (a == "organization_administration" && b == "members") || (a == "members" && b == "organization_administration"):
		return "Organization administration and Members write together enable promoting attacker accounts to Org Owner and locking out legitimate admins."
	default:
		return ""
	}
}

func writeupExploitPath(section, a, b string) string {
	if (a == "checks" && b == "pull_requests") || (a == "pull_requests" && b == "checks") {
		return "Create a pull request, mark required checks as passed via the Checks API, and merge without legitimate CI or human approval."
	}
	lower := strings.ToLower(section)
	if strings.Contains(lower, "promote attacker account to org owner") {
		return "Promote an attacker-controlled account to Org Owner, then demote or remove legitimate admins."
	}
	return ""
}
