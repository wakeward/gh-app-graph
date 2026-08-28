package inventory

import (
	"encoding/json"
	"testing"
	"time"
)

const fixtureJSON = `{
  "paths": {
    "/repos/{owner}/{repo}/contents/{path}": {
      "GET": {"access": "read", "permission": "contents"},
      "PUT": {"access": "write", "permission": "contents"}
    },
    "/repos/{owner}/{repo}/actions/workflows/{workflow_id}": {
      "GET": {"access": "read", "permission": "actions"}
    }
  },
  "permissions": {
    "contents": {
      "read": ["GET /repos/{owner}/{repo}/contents/{path}"],
      "write": ["PUT /repos/{owner}/{repo}/contents/{path}"],
      "url": "https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps#contents"
    },
    "single_file": {
      "read": ["GET /repos/{owner}/{repo}/contents/{path}"],
      "write": ["PUT /repos/{owner}/{repo}/contents/{path}"],
      "url": "https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps#single-file"
    },
    "actions": {
      "read": ["GET /repos/{owner}/{repo}/actions/workflows/{workflow_id}"],
      "url": "https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps#actions"
    }
  }
}`

func TestNormalize_BuildsEndpointIndexOverlap(t *testing.T) {
	inv, err := Normalize("file://fixture", []byte(fixtureJSON), time.Date(2026, 7, 28, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("Normalize: %v", err)
	}

	key := "GET /repos/{owner}/{repo}/contents/{path}"
	grants := inv.EndpointIndex[key]
	if len(grants) != 2 {
		t.Fatalf("expected 2 permissions for %q, got %d: %+v", key, len(grants), grants)
	}

	multi := 0
	for _, g := range inv.EndpointIndex {
		if len(g) > 1 {
			multi++
		}
	}
	if multi != 2 {
		t.Errorf("expected 2 multi-permission endpoints (GET+PUT contents path), got %d", multi)
	}
	if inv.Meta.GitHubDocsAPIVersion != GitHubDocsAPIVersion {
		t.Errorf("meta api version = %q, want %q", inv.Meta.GitHubDocsAPIVersion, GitHubDocsAPIVersion)
	}
	if inv.Meta.EndpointCount != 3 {
		t.Errorf("endpoint count = %d, want 3", inv.Meta.EndpointCount)
	}
	if inv.Meta.PermissionCount != 3 {
		t.Errorf("permission count = %d, want 3", inv.Meta.PermissionCount)
	}
}

func TestNormalize_AdminAccessMapsToWrite(t *testing.T) {
	raw := `{
	  "paths": {"/orgs/{org}/properties/schema": {"PATCH": {"access": "admin", "permission": "organization_custom_properties"}}},
	  "permissions": {"organization_custom_properties": {"write": ["PATCH /orgs/{org}/properties/schema"], "url": "https://example.com"}}
	}`
	inv, err := Normalize("test", []byte(raw), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if inv.Endpoints[0].Access != "write" {
		t.Errorf("admin access mapped to %q, want write", inv.Endpoints[0].Access)
	}
}

func TestWriteJSON_RoundTrip(t *testing.T) {
	inv, err := Normalize("test", []byte(fixtureJSON), time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.MarshalIndent(inv, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	var decoded Inventory
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if len(decoded.EndpointIndex) != len(inv.EndpointIndex) {
		t.Fatalf("round-trip endpoint_index length mismatch")
	}
}
