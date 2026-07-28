package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/wakeward/github-app-permissions-graph/pkg/model"
)

type csvRow struct {
	Technique   string
	Combination string
	Risk        string
	ExploitPath string
}

var (
	grantSegmentRE = regexp.MustCompile(`^\s*([^:(]+?)\s*:\s*(read|write)\s*$`)
	orSuffixRE     = regexp.MustCompile(`(?i)^(.+?)\s*\(or\s+(.+)\)\s*$`)
)

func readCSV(path string) ([]csvRow, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("empty CSV")
	}

	header := records[0]
	idx := make(map[string]int, len(header))
	for i, h := range header {
		idx[strings.TrimSpace(h)] = i
	}
	required := []string{"Technique", "Toxic Combination", "Risk", "Exploit Path"}
	for _, col := range required {
		if _, ok := idx[col]; !ok {
			return nil, fmt.Errorf("missing expected column %q", col)
		}
	}

	get := func(rec []string, col string) string {
		i := idx[col]
		if i >= len(rec) {
			return ""
		}
		return strings.TrimSpace(rec[i])
	}

	var rows []csvRow
	for _, rec := range records[1:] {
		technique := get(rec, "Technique")
		if technique == "" {
			continue
		}
		rows = append(rows, csvRow{
			Technique:   technique,
			Combination: get(rec, "Toxic Combination"),
			Risk:        get(rec, "Risk"),
			ExploitPath: get(rec, "Exploit Path"),
		})
	}
	return rows, nil
}

// parseCombination expands a CSV "Toxic Combination" cell into one or more
// grant sets. Alternatives written as "(or Other: write)" become separate
// combinations with the same technique metadata.
func parseCombination(text string, known map[string]struct{}) ([][]model.PermissionGrant, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("empty combination")
	}

	// Single segment that is only an OR choice, e.g.
	// "Organization administration: write (or Members: write)".
	if !strings.Contains(text, "+") {
		if m := orSuffixRE.FindStringSubmatch(text); len(m) == 3 {
			left, err := parseGrantSegment(m[1], known)
			if err != nil {
				return nil, err
			}
			right, err := parseGrantSegment(m[2], known)
			if err != nil {
				return nil, err
			}
			return [][]model.PermissionGrant{{left}, {right}}, nil
		}
		grant, err := parseGrantSegment(text, known)
		if err != nil {
			return nil, err
		}
		return [][]model.PermissionGrant{{grant}}, nil
	}

	parts := strings.Split(text, "+")
	base := make([]model.PermissionGrant, 0, len(parts))
	var variants [][]model.PermissionGrant

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if m := orSuffixRE.FindStringSubmatch(part); len(m) == 3 {
			left, err := parseGrantSegment(m[1], known)
			if err != nil {
				return nil, err
			}
			right, err := parseGrantSegment(m[2], known)
			if err != nil {
				return nil, err
			}
			variants = append(variants, append(append([]model.PermissionGrant{}, base...), left))
			variants = append(variants, append(append([]model.PermissionGrant{}, base...), right))
			continue
		}
		grant, err := parseGrantSegment(part, known)
		if err != nil {
			return nil, err
		}
		base = append(base, grant)
	}

	if len(variants) > 0 {
		return variants, nil
	}
	return [][]model.PermissionGrant{base}, nil
}

func parseGrantSegment(segment string, known map[string]struct{}) (model.PermissionGrant, error) {
	m := grantSegmentRE.FindStringSubmatch(strings.TrimSpace(segment))
	if len(m) != 3 {
		return model.PermissionGrant{}, fmt.Errorf("cannot parse grant segment %q", segment)
	}
	key, err := resolveAPIKey(m[1], known)
	if err != nil {
		return model.PermissionGrant{}, err
	}
	access, err := parseAccessLevel(m[2])
	if err != nil {
		return model.PermissionGrant{}, err
	}
	return model.PermissionGrant{APIKey: key, Access: access}, nil
}

func inferBlastRadius(technique, risk string) model.BlastRadius {
	text := strings.ToLower(technique + " " + risk)
	switch {
	case strings.Contains(text, "catastrophic"), strings.Contains(text, "critical"), strings.Contains(text, "complete organization"), strings.Contains(text, "complete loss"):
		return model.BlastRadiusCritical
	case strings.Contains(text, "high to critical"):
		return model.BlastRadiusCritical
	case strings.Contains(text, "high"), strings.Contains(text, "enterprise-wide"):
		return model.BlastRadiusHigh
	case strings.Contains(text, "medium"):
		return model.BlastRadiusMedium
	}
	switch slugify(technique) {
	case "stealth-backdoor", "arbitrary-code-execution", "organization-takeover":
		return model.BlastRadiusCritical
	default:
		return model.BlastRadiusHigh
	}
}

func buildFromCSV(rows []csvRow, known map[string]struct{}) ([]model.ToxicCombination, []string, error) {
	var out []model.ToxicCombination
	var warnings []string

	for _, row := range rows {
		grantSets, err := parseCombination(row.Combination, known)
		if err != nil {
			return nil, warnings, fmt.Errorf("%s: %w", row.Technique, err)
		}
		baseID := slugify(row.Technique)
		blast := inferBlastRadius(row.Technique, row.Risk)

		for i, grants := range grantSets {
			if len(grants) < 2 {
				warnings = append(warnings, fmt.Sprintf("%s: skipped single-permission variant %s (need 2+ grants for a toxic combination)", row.Technique, grantSetKey(grants)))
				continue
			}
		id := baseID
		if len(grantSets) > 1 {
			suffix := variantSuffix(grants)
			if suffix != "" {
				id = baseID + "-" + suffix
			} else {
				id = baseID + "-" + slugify(grantSetKey(grants))
			}
		} else if i > 0 {
				id = fmt.Sprintf("%s-%d", baseID, i+1)
			}
			out = append(out, model.ToxicCombination{
				ID:              id,
				Technique:       row.Technique,
				Permissions:     grants,
				BlastRadius:     blast,
				RiskDescription: row.Risk,
				ExploitPath:     row.ExploitPath,
				Source:          "migrate-toxic-combinations CSV",
			})
		}
	}
	return out, warnings, nil
}

func variantSuffix(grants []model.PermissionGrant) string {
	keys := make(map[string]struct{}, len(grants))
	for _, g := range grants {
		keys[g.APIKey] = struct{}{}
	}
	_, hasContents := keys["contents"]
	_, hasDeployments := keys["deployments"]
	_, hasPackages := keys["packages"]
	switch {
	case hasContents && hasDeployments:
		return "deployments"
	case hasContents && hasPackages:
		return "packages"
	default:
		return ""
	}
}
