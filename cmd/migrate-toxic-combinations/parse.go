package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/wakeward/gh-app-graph/pkg/model"
)

type csvRow struct {
	Technique   string
	Combination string
	Risk        string
	ExploitPath string
}

var (
	grantSegmentRE = regexp.MustCompile(`^\s*([^:(]+?)\s*:\s*(read|write)\s*$`)
	orSuffixRE     = regexp.MustCompile(`(?i)^(.+?)\s*\(or\s+(.+?)\)?\s*$`)
)

func readCSV(path string) ([]csvRow, error) {
	f, err := os.Open(path) // #nosec G304 -- path is a caller-supplied CSV for one-time migration
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
	required := []string{"Technique", "Toxic Combination", "Exploit Path"}
	for _, col := range required {
		if _, ok := idx[col]; !ok {
			return nil, fmt.Errorf("missing expected column %q", col)
		}
	}
	severityCol := "Severity"
	if _, ok := idx[severityCol]; !ok {
		severityCol = "Risk"
		if _, ok := idx[severityCol]; !ok {
			return nil, fmt.Errorf("missing expected column %q or %q", "Severity", "Risk")
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
			Risk:        get(rec, severityCol),
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

	if !strings.Contains(text, "+") {
		opts, err := expandPartOptions(text, known)
		if err != nil {
			return nil, err
		}
		sets := make([][]model.PermissionGrant, len(opts))
		for i, opt := range opts {
			sets[i] = []model.PermissionGrant{opt}
		}
		return sets, nil
	}

	parts := strings.Split(text, "+")
	bases := [][]model.PermissionGrant{{}}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		part = strings.TrimPrefix(part, "(")
		part = strings.TrimSuffix(part, ")")
		part = strings.TrimSpace(part)

		opts, err := expandPartOptions(part, known)
		if err != nil {
			return nil, err
		}

		var next [][]model.PermissionGrant
		for _, base := range bases {
			for _, opt := range opts {
				next = append(next, append(append([]model.PermissionGrant{}, base...), opt))
			}
		}
		bases = next
	}
	return bases, nil
}

// expandPartOptions returns one or more grant alternatives for a single
// plus-separated segment (handles parenthesized OR and bare "or" joins).
func expandPartOptions(part string, known map[string]struct{}) ([]model.PermissionGrant, error) {
	part = strings.TrimSpace(part)
	if part == "" {
		return nil, fmt.Errorf("empty combination segment")
	}

	if m := orSuffixRE.FindStringSubmatch(part); len(m) == 3 {
		left, right, err := parseOrAlternatives(m[1], m[2], known)
		if err != nil {
			return nil, err
		}
		return []model.PermissionGrant{left, right}, nil
	}

	if strings.Contains(strings.ToLower(part), " or ") {
		alts := strings.Split(part, " or ")
		out := make([]model.PermissionGrant, 0, len(alts))
		for _, alt := range alts {
			grant, err := parseGrantSegment(strings.TrimSpace(alt), known)
			if err != nil {
				return nil, err
			}
			out = append(out, grant)
		}
		return out, nil
	}

	grant, err := parseGrantSegment(part, known)
	if err != nil {
		return nil, err
	}
	return []model.PermissionGrant{grant}, nil
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

// parseOrAlternatives parses "Left: write (or right)" segments. When the
// right-hand alternative omits an access level (e.g. bare api_key), it
// inherits the left-hand access.
func parseOrAlternatives(leftText, rightText string, known map[string]struct{}) (model.PermissionGrant, model.PermissionGrant, error) {
	left, err := parseGrantSegment(leftText, known)
	if err != nil {
		return model.PermissionGrant{}, model.PermissionGrant{}, err
	}
	right, err := parseGrantSegment(rightText, known)
	if err == nil {
		return left, right, nil
	}
	key, err := resolveAPIKey(strings.TrimSpace(rightText), known)
	if err != nil {
		return model.PermissionGrant{}, model.PermissionGrant{}, fmt.Errorf("cannot parse grant segment %q", rightText)
	}
	return left, model.PermissionGrant{APIKey: key, Access: left.Access}, nil
}

func normalizeSeverityWord(w string) (model.BlastRadius, bool) {
	switch strings.ToLower(strings.TrimSpace(w)) {
	case "critical":
		return model.BlastRadiusCritical, true
	case "high":
		return model.BlastRadiusHigh, true
	case "medium":
		return model.BlastRadiusMedium, true
	case "low":
		return model.BlastRadiusMedium, true
	default:
		return "", false
	}
}

func riskDescription(row csvRow) string {
	if _, ok := normalizeSeverityWord(row.Risk); ok {
		return fmt.Sprintf("%s (%s)", row.Technique, row.Risk)
	}
	return row.Risk
}

func inferComboPlatform(combo model.ToxicCombination, ghesKeys map[string]struct{}) model.PlatformAvailability {
	if strings.Contains(strings.ToLower(combo.ExploitPath), "ghes only") {
		return model.PlatformGHESOnly
	}
	for _, grant := range combo.Permissions {
		if _, ok := ghesKeys[grant.APIKey]; ok {
			return model.PlatformGHESOnly
		}
	}
	return model.PlatformAll
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
				RiskDescription: riskDescription(row),
				ExploitPath:     row.ExploitPath,
				Source:          "migrate-toxic-combinations CSV",
			})
		}
	}
	return out, warnings, nil
}

func variantSuffix(grants []model.PermissionGrant) string {
	if len(grants) == 1 {
		return grants[0].APIKey
	}
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
