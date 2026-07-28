package main

import (
	"testing"

	"github.com/wakeward/github-app-permissions-graph/pkg/model"
)

func testKnownKeys() map[string]struct{} {
	keys := []string{
		"administration", "contents", "workflows", "organization_administration",
		"members", "deployments", "packages", "checks", "pull_requests",
	}
	out := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		out[k] = struct{}{}
	}
	return out
}

func TestParseCombination_PlusSeparated(t *testing.T) {
	known := testKnownKeys()
	sets, err := parseCombination("Administration: write + Contents: write", known)
	if err != nil {
		t.Fatal(err)
	}
	if len(sets) != 1 || len(sets[0]) != 2 {
		t.Fatalf("expected one 2-grant set, got %+v", sets)
	}
	if sets[0][0].APIKey != "administration" || sets[0][1].APIKey != "contents" {
		t.Fatalf("unexpected grants: %+v", sets[0])
	}
}

func TestParseCombination_OrAlternatives(t *testing.T) {
	known := testKnownKeys()
	sets, err := parseCombination("Contents: write + Deployments: write (or Packages: write)", known)
	if err != nil {
		t.Fatal(err)
	}
	if len(sets) != 2 {
		t.Fatalf("expected 2 variants, got %d", len(sets))
	}
	if grantSetKey(sets[0]) != "contents:write+deployments:write" {
		t.Errorf("variant 0 = %s", grantSetKey(sets[0]))
	}
	if grantSetKey(sets[1]) != "contents:write+packages:write" {
		t.Errorf("variant 1 = %s", grantSetKey(sets[1]))
	}
}

func TestParseCombination_SingleOrSegmentSkippedInBuild(t *testing.T) {
	known := testKnownKeys()
	rows := []csvRow{{
		Technique:   "Organization Takeover",
		Combination: "Organization administration: write (or Members: write)",
		Risk:        "Org-level catastrophic blast radius.",
		ExploitPath: "Invite malicious owners.",
	}}
	combos, warnings, err := buildFromCSV(rows, known)
	if err != nil {
		t.Fatal(err)
	}
	if len(combos) != 0 {
		t.Fatalf("expected single-permission variants to be skipped, got %d combos", len(combos))
	}
	if len(warnings) != 2 {
		t.Fatalf("expected 2 skip warnings, got %d: %v", len(warnings), warnings)
	}
}

func TestInferBlastRadius(t *testing.T) {
	if got := inferBlastRadius("Supply Chain Poisoning via Releases", "Allows manipulation of release artifacts"); got != model.BlastRadiusHigh {
		t.Errorf("got %q, want High", got)
	}
	if got := inferBlastRadius("Stealth Backdoor", "Administration: write allows..."); got != model.BlastRadiusCritical {
		t.Errorf("got %q, want Critical", got)
	}
}

func TestTechniqueForWriteupPair(t *testing.T) {
	pair := grantSetKey([]model.PermissionGrant{
		{APIKey: "organization_administration", Access: model.AccessWrite},
		{APIKey: "members", Access: model.AccessWrite},
	})
	t.Logf("pair=%q", pair)
	if got := techniqueForWriteupPair("organization_administration", "members"); got != "Organization Takeover" {
		t.Errorf("got %q, want Organization Takeover", got)
	}
}

func TestVariantSuffix(t *testing.T) {
	grants := []model.PermissionGrant{
		{APIKey: "contents", Access: model.AccessWrite},
		{APIKey: "packages", Access: model.AccessWrite},
	}
	if got := variantSuffix(grants); got != "packages" {
		t.Errorf("got %q, want packages", got)
	}
}

func TestResolveAPIKey_OrganizationAdministration(t *testing.T) {
	known := testKnownKeys()
	key, err := resolveAPIKey("Organization administration", known)
	if err != nil || key != "organization_administration" {
		t.Fatalf("resolveAPIKey() = (%q, %v)", key, err)
	}
}
