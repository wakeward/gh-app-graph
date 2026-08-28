// Command migrate-csv is a one-time migration: it converts the original
// "GitHub Apps Permissions" CSV into data/permissions/*.yaml and seeds the
// undocumented Enterprise permissions visible only in the App creation UI
// screenshot. It does not archive the source CSV/screenshot into the repo -
// those are working documents, not permanent provenance artifacts.
//
// Usage:
//
//	go run ./cmd/migrate-csv --csv "/path/to/GitHub Apps Permissions - GitHub-Apps-Permissions.csv"
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wakeward/gh-app-graph/pkg/fileio"
	"github.com/wakeward/gh-app-graph/pkg/model"
)

func main() {
	csvPath := flag.String("csv", "", "path to the source GitHub Apps Permissions CSV (required)")
	outDir := flag.String("out", "data/permissions", "output directory for the generated *.yaml files")
	flag.Parse()

	if *csvPath == "" {
		fmt.Fprintln(os.Stderr, "migrate-csv: --csv is required")
		os.Exit(1)
	}

	rows, err := readCSV(*csvPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate-csv: reading CSV: %v\n", err)
		os.Exit(1)
	}

	permissions, warnings, err := buildPermissions(rows)
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate-csv: %v\n", err)
		os.Exit(1)
	}
	permissions = append(permissions, seedUndocumentedEnterprisePermissions()...)

	if err := writeByCategory(permissions, *outDir); err != nil {
		fmt.Fprintf(os.Stderr, "migrate-csv: writing output: %v\n", err)
		os.Exit(1)
	}

	needsInvestigation := 0
	for _, p := range permissions {
		if p.NeedsInvestigation {
			needsInvestigation++
		}
	}
	fmt.Printf("migrate-csv: wrote %d permissions (%d flagged needs_investigation) across %s\n", len(permissions), needsInvestigation, *outDir)
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "migrate-csv: warning: %s\n", w)
	}
}

// csvRow is one raw row of the source CSV.
type csvRow struct {
	Group         string
	Name          string
	APIKey        string
	AccessText    string
	Overview      string
	SeverityCol   string
	Notes         string
	FurtherInvest string
}

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
	required := []string{"Permission Group", "Permission Name", "API Key", "Access Level", "Overview", "Severity", "Further Investigation?"}
	for _, col := range required {
		if _, ok := idx[col]; !ok {
			return nil, fmt.Errorf("missing expected column %q", col)
		}
	}
	notesCol := "Security Context"
	if _, ok := idx[notesCol]; !ok {
		notesCol = "Initial Security Thoughts"
		if _, ok := idx[notesCol]; !ok {
			return nil, fmt.Errorf("missing expected column %q or %q", "Security Context", "Initial Security Thoughts")
		}
	}

	get := func(rec []string, col string) string {
		i := idx[col]
		if i >= len(rec) {
			return ""
		}
		return rec[i]
	}

	var rows []csvRow
	for _, rec := range records[1:] {
		apiKey := strings.TrimSpace(get(rec, "API Key"))
		if apiKey == "" {
			continue
		}
		rows = append(rows, csvRow{
			Group:         strings.TrimSpace(get(rec, "Permission Group")),
			Name:          strings.TrimSpace(get(rec, "Permission Name")),
			APIKey:        apiKey,
			AccessText:    strings.TrimSpace(get(rec, "Access Level")),
			Overview:      strings.TrimSpace(get(rec, "Overview")),
			SeverityCol:   strings.TrimSpace(get(rec, "Severity")),
			Notes:         strings.TrimSpace(get(rec, notesCol)),
			FurtherInvest: strings.TrimSpace(get(rec, "Further Investigation?")),
		})
	}
	return rows, nil
}

func buildPermissions(rows []csvRow) ([]model.Permission, []string, error) {
	type group struct {
		group, name, overview string
		levels                []model.AccessLevelDetail
		needsInvestigation    bool
	}
	order := make([]string, 0)
	groups := make(map[string]*group)
	var warnings []string

	for _, row := range rows {
		g, ok := groups[row.APIKey]
		if !ok {
			g = &group{group: row.Group, name: row.Name}
			groups[row.APIKey] = g
			order = append(order, row.APIKey)
		}
		if len(row.Overview) > len(g.overview) {
			g.overview = row.Overview
		}

		access, err := parseAccess(row.AccessText)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s: %v (row skipped)", row.APIKey, err))
			continue
		}

		severity, notes, ok := parseSeverity(row.SeverityCol, row.Notes)
		if !ok {
			severity = model.SeverityUnknown
			g.needsInvestigation = true
			warnings = append(warnings, fmt.Sprintf("%s (%s): could not parse severity from %q / %q - marked Unknown, needs_investigation", row.APIKey, access, row.SeverityCol, row.Notes))
		}
		if strings.EqualFold(row.FurtherInvest, "yes") {
			g.needsInvestigation = true
		}

		g.levels = append(g.levels, model.AccessLevelDetail{
			Access:        access,
			Severity:      severity,
			SecurityNotes: notes,
		})
	}

	permissions := make([]model.Permission, 0, len(order))
	for _, key := range order {
		g := groups[key]
		plane, ok := classifyImpactPlane(key)
		if !ok {
			g.needsInvestigation = true
			warnings = append(warnings, fmt.Sprintf("%s: no impact_plane classification - defaulting to data_execution, needs_investigation set", key))
			plane = model.ImpactPlaneDataExecution
		}
		permissions = append(permissions, model.Permission{
			Name:                 g.name,
			APIKey:               key,
			Category:             strings.ToLower(g.group),
			Overview:             g.overview,
			AccessLevels:         g.levels,
			NeedsInvestigation:   g.needsInvestigation,
			DocStatus:            inferDocStatus(g.overview),
			ImpactPlane:          plane,
			PlatformAvailability: inferPlatformAvailability(g.overview, g.levels),
		})
	}
	return permissions, warnings, nil
}

func parseAccess(s string) (model.AccessLevel, error) {
	switch s {
	case "Read / Write":
		return model.AccessWrite, nil
	case "Read-Only":
		return model.AccessRead, nil
	default:
		return "", fmt.Errorf("unrecognized access level %q", s)
	}
}

// normalizeSeverityWord maps a bare severity word (from either the Severity
// column or the leading word of the notes column) to a model.Severity.
func normalizeSeverityWord(w string) (model.Severity, bool) {
	switch strings.ToUpper(strings.TrimSpace(w)) {
	case "CRITICAL":
		return model.SeverityCritical, true
	case "HIGH":
		return model.SeverityHigh, true
	case "MEDIUM":
		return model.SeverityMedium, true
	case "LOW":
		return model.SeverityLow, true
	case "N/A", "NA":
		return model.SeverityInformational, true
	default:
		return "", false
	}
}

// parseSeverity determines a row's severity from the Severity column if
// present, otherwise from a leading severity word embedded in the notes
// column (e.g. "HIGH. Could expand..."), stripping that word from the notes
// returned so it isn't duplicated. Returns ok=false if neither source
// yields a recognized severity.
func parseSeverity(severityCol, notes string) (model.Severity, string, bool) {
	if s := strings.TrimSpace(severityCol); s != "" {
		if sev, ok := normalizeSeverityWord(s); ok {
			return sev, strings.TrimSpace(notes), true
		}
	}

	trimmed := strings.TrimSpace(notes)
	if trimmed == "" {
		return "", "", false
	}
	word := trimmed
	rest := ""
	if idx := strings.IndexAny(trimmed, ". "); idx > 0 {
		word = trimmed[:idx]
		rest = strings.TrimSpace(strings.TrimPrefix(trimmed[idx:], "."))
	}
	if sev, ok := normalizeSeverityWord(word); ok {
		return sev, rest, true
	}
	return "", trimmed, false
}

func inferPlatformAvailability(overview string, levels []model.AccessLevelDetail) model.PlatformAvailability {
	text := strings.ToLower(overview)
	for _, level := range levels {
		text += " " + strings.ToLower(level.SecurityNotes)
	}
	if strings.Contains(text, "github enterprise server") || strings.Contains(text, "ghes only") {
		return model.PlatformGHESOnly
	}
	return model.PlatformAll
}

func inferDocStatus(overview string) model.DocStatus {
	lower := strings.ToLower(overview)
	switch {
	case strings.Contains(overview, "UNCONFIRMED KEY"):
		return model.DocStatusUnconfirmedKey
	case strings.Contains(lower, "not found in"):
		return model.DocStatusDisputed
	case strings.Contains(lower, "not observed in the live"):
		return model.DocStatusDisputed
	case strings.Contains(overview, "(Preview)"):
		return model.DocStatusUndocumentedPreview
	default:
		return model.DocStatusDocumented
	}
}

func writeByCategory(permissions []model.Permission, outDir string) error {
	byCategory := make(map[string][]model.Permission)
	for _, p := range permissions {
		byCategory[p.Category] = append(byCategory[p.Category], p)
	}

	if err := fileio.MkdirAll(outDir); err != nil {
		return err
	}

	for category, perms := range byCategory {
		sort.Slice(perms, func(i, j int) bool { return perms[i].APIKey < perms[j].APIKey })

		header := fmt.Sprintf(`# DRAFT %s permissions - NOT CANONICAL
# Seeded by cmd/migrate-csv from an incomplete working spreadsheet plus Enterprise
# UI observations. Severity/impact_plane are starting notes only. Replace or reconcile
# once fetch-inventory and detect-overlap produce the real endpoint graph.
# See data/permissions/README.md
`, strings.Title(category))
		outPath := filepath.Join(outDir, category+".yaml")
		if err := fileio.WriteYAML(outPath, header, perms); err != nil { // #nosec G117 -- api_key fields are permission identifiers
			return err
		}
	}
	return nil
}
