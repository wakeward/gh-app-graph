package main

import (
	"testing"

	"github.com/wakeward/gh-app-graph/pkg/model"
)

func TestParseSeverity_FromColumn(t *testing.T) {
	sev, notes, ok := parseSeverity("Low", "Nuisance vector")
	if !ok || sev != model.SeverityLow || notes != "Nuisance vector" {
		t.Errorf("got (%q, %q, %v), want (Low, %q, true)", sev, notes, ok, "Nuisance vector")
	}
}

func TestParseSeverity_ColumnNATreatedAsInformational(t *testing.T) {
	sev, _, ok := parseSeverity("N/A", "")
	if !ok || sev != model.SeverityInformational {
		t.Errorf("got (%q, %v), want (Informational, true)", sev, ok)
	}
}

func TestParseSeverity_FromNotesLeadingWord(t *testing.T) {
	sev, notes, ok := parseSeverity("", "HIGH. Could expand the blast radius.")
	if !ok || sev != model.SeverityHigh {
		t.Fatalf("got (%q, ok=%v), want (High, true)", sev, ok)
	}
	if notes != "Could expand the blast radius." {
		t.Errorf("expected severity prefix stripped from notes, got %q", notes)
	}
}

func TestParseSeverity_UnparseableFallsThrough(t *testing.T) {
	sev, notes, ok := parseSeverity("", "What does RO allow? [FI]")
	if ok {
		t.Fatalf("expected ok=false for unparseable notes, got severity=%q", sev)
	}
	if notes != "What does RO allow? [FI]" {
		t.Errorf("expected original notes preserved on failure, got %q", notes)
	}
}

func TestParseSeverity_EmptyBoth(t *testing.T) {
	_, _, ok := parseSeverity("", "")
	if ok {
		t.Fatalf("expected ok=false when both column and notes are empty")
	}
}

func TestParseAccess(t *testing.T) {
	cases := map[string]model.AccessLevel{
		"Read / Write": model.AccessWrite,
		"Read-Only":    model.AccessRead,
	}
	for text, want := range cases {
		got, err := parseAccess(text)
		if err != nil {
			t.Errorf("parseAccess(%q) unexpected error: %v", text, err)
		}
		if got != want {
			t.Errorf("parseAccess(%q) = %q, want %q", text, got, want)
		}
	}
	if _, err := parseAccess("Something Else"); err == nil {
		t.Errorf("expected error for unrecognized access text")
	}
}

func TestInferDocStatus(t *testing.T) {
	cases := []struct {
		overview string
		want     model.DocStatus
	}{
		{"Manage enterprise custom properties. (Preview)", model.DocStatusUndocumentedPreview},
		{"UNCONFIRMED KEY - inferred from an endpoint path.", model.DocStatusUnconfirmedKey},
		{"NOT FOUND in go-github's schema and NOT FOUND in the live UI.", model.DocStatusDisputed},
		{"Confirmed as a field in go-github's schema, but NOT observed in the live GitHub App permissions UI.", model.DocStatusDisputed},
		{"Access repository contents, commits, and branches.", model.DocStatusDocumented},
	}
	for _, c := range cases {
		if got := inferDocStatus(c.overview); got != c.want {
			t.Errorf("inferDocStatus(%q) = %q, want %q", c.overview, got, c.want)
		}
	}
}

func TestClassifyImpactPlane_EverySeedEntryIsClassified(t *testing.T) {
	for _, p := range seedUndocumentedEnterprisePermissions() {
		if p.ImpactPlane == "" {
			t.Errorf("seed entry %s has no impact_plane", p.APIKey)
		}
	}
}

func TestClassifyImpactPlane_MissingKeyReturnsNotOK(t *testing.T) {
	if _, ok := classifyImpactPlane("definitely_not_a_real_key"); ok {
		t.Errorf("expected ok=false for an unclassified key")
	}
}
