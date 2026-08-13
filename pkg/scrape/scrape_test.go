package scrape

import (
	"testing"

	"github.com/wakeward/github-app-permissions-graph/pkg/inventory"
	"github.com/wakeward/github-app-permissions-graph/pkg/model"
)

func TestSelectTargets_DefaultHighCriticalOnly(t *testing.T) {
	perms := []model.Permission{{
		APIKey: "contents",
		AccessLevels: []model.AccessLevelDetail{
			{Access: model.AccessRead, Severity: model.SeverityLow},
			{Access: model.AccessWrite, Severity: model.SeverityCritical},
		},
	}}
	targets := SelectTargets(perms, map[string]string{"contents": "https://example.com/contents"}, false)
	if len(targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(targets))
	}
	if targets[0].Access != model.AccessWrite {
		t.Fatalf("expected write access target, got %s", targets[0].Access)
	}
}

func TestSelectTargets_WidenNeedsInvestigation(t *testing.T) {
	perms := []model.Permission{{
		APIKey:             "keys",
		NeedsInvestigation: true,
		AccessLevels: []model.AccessLevelDetail{
			{Access: model.AccessRead, Severity: model.SeverityUnknown},
		},
	}}
	targets := SelectTargets(perms, nil, true)
	if len(targets) != 1 {
		t.Fatalf("expected widen to include needs_investigation row, got %d", len(targets))
	}
}

func TestExtractRequirements_FindsGrantAndExcludesSource(t *testing.T) {
	text := "This endpoint also requires the `contents:write` permission and the `metadata` permission."
	got := ExtractRequirements(text, "workflows")
	if len(got) != 2 {
		t.Fatalf("expected 2 requirements, got %d: %+v", len(got), got)
	}
}

func TestExtractFromHTML_StripsTags(t *testing.T) {
	html := `<html><body><p>Requires the <code>contents</code> permission</p></body></html>`
	text := ExtractFromHTML(html)
	if text == "" || text == html {
		t.Fatalf("expected stripped text, got %q", text)
	}
}

func TestDetectProseDependencies_UsesFetcher(t *testing.T) {
	perms := []model.Permission{{
		APIKey: "workflows",
		AccessLevels: []model.AccessLevelDetail{
			{Access: model.AccessWrite, Severity: model.SeverityCritical},
		},
	}}
	inv := &inventory.Inventory{
		Permissions: map[string]inventory.PermissionEndpoints{
			"workflows": {DocURL: "https://example.com/workflows"},
		},
	}
	fetcher := func(url string) (string, error) {
		return `<p>Requires the ` + "`contents:write`" + ` permission for workflow uploads.</p>`, nil
	}
	edges, warnings, err := DetectProseDependencies(perms, inv, false, fetcher)
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(edges) != 1 || len(edges[0].Requires) != 1 {
		t.Fatalf("expected one edge with one require, got %+v", edges)
	}
	if edges[0].Requires[0].Permission != "contents" {
		t.Fatalf("got %+v", edges[0].Requires[0])
	}
}
