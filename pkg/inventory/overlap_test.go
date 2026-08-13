package inventory

import (
	"testing"

	"github.com/wakeward/gh-app-graph/pkg/model"
)

func TestDetectOverlapDependencies_SingleFileRequiresContents(t *testing.T) {
	inv := &Inventory{
		Permissions: map[string]PermissionEndpoints{
			"contents":    {DocURL: "https://docs.github.com/contents"},
			"single_file": {DocURL: "https://docs.github.com/single-file"},
		},
		EndpointIndex: map[string][]PermissionGrant{
			"GET /repos/{owner}/{repo}/contents/{path}": {
				{Permission: "contents", Access: "read"},
				{Permission: "single_file", Access: "read"},
			},
			"PUT /repos/{owner}/{repo}/contents/{path}": {
				{Permission: "contents", Access: "write"},
				{Permission: "single_file", Access: "write"},
			},
		},
	}

	edges := DetectOverlapDependencies(inv)

	var singleFileWrite *model.DependencyEdge
	for i := range edges {
		if edges[i].Permission == "single_file" && edges[i].Access == model.AccessWrite {
			singleFileWrite = &edges[i]
			break
		}
	}
	if singleFileWrite == nil {
		t.Fatalf("expected single_file:write edge, got %+v", edges)
	}
	if len(singleFileWrite.Requires) != 1 || singleFileWrite.Requires[0].Permission != "contents" {
		t.Fatalf("expected single_file:write requires contents:write, got %+v", singleFileWrite.Requires)
	}
	if singleFileWrite.Requires[0].Type != model.DependencyTypeCategoryOverlap {
		t.Errorf("expected category-overlap type")
	}
}

func TestDetectOverlapDependencies_SkipsSinglePermissionEndpoints(t *testing.T) {
	inv := &Inventory{
		EndpointIndex: map[string][]PermissionGrant{
			"GET /meta": {{Permission: "metadata", Access: "read"}},
		},
	}
	if edges := DetectOverlapDependencies(inv); len(edges) != 0 {
		t.Fatalf("expected no edges, got %d", len(edges))
	}
}

func TestDetectOverlapDependencies_DeterministicOrder(t *testing.T) {
	inv := &Inventory{
		Permissions: map[string]PermissionEndpoints{
			"a": {}, "b": {}, "c": {},
		},
		EndpointIndex: map[string][]PermissionGrant{
			"GET /x": {
				{Permission: "c", Access: "read"},
				{Permission: "a", Access: "read"},
				{Permission: "b", Access: "read"},
			},
		},
	}
	e1 := DetectOverlapDependencies(inv)
	e2 := DetectOverlapDependencies(inv)
	if len(e1) != len(e2) {
		t.Fatal("non-deterministic edge count")
	}
	for i := range e1 {
		if e1[i].Permission != e2[i].Permission || e1[i].Access != e2[i].Access {
			t.Fatalf("edge order differed at %d", i)
		}
	}
}
