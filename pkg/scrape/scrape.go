package scrape

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/wakeward/gh-app-graph/pkg/inventory"
	"github.com/wakeward/gh-app-graph/pkg/model"
)

// Fetcher downloads a docs page. Production code uses HTTP; tests inject fixtures.
type Fetcher func(url string) (string, error)

// DefaultFetcher fetches HTML from GitHub docs with a short timeout.
func DefaultFetcher(url string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "gh-app-graph/scrape-prose")
	resp, err := client.Do(req) // #nosec G107 -- docs URLs come from inventory doc_url fields
	if err != nil {
		return "", fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s: HTTP %s", url, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// DetectProseDependencies scrapes permission doc pages for Type B draft edges.
func DetectProseDependencies(permissions []model.Permission, inv *inventory.Inventory, widen bool, fetch Fetcher) ([]model.DependencyEdge, []string, error) {
	if fetch == nil {
		fetch = DefaultFetcher
	}
	docURLs := docURLsFromInventory(inv)
	targets := SelectTargets(permissions, docURLs, widen)

	edgeByKey := make(map[string]*model.DependencyEdge)
	var warnings []string

	for _, target := range targets {
		if target.DocURL == "" {
			warnings = append(warnings, fmt.Sprintf("%s:%s: no doc_url in inventory; skipped", target.Permission.APIKey, target.Access))
			continue
		}
		html, err := fetch(target.DocURL)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("%s:%s: fetch %s: %v", target.Permission.APIKey, target.Access, target.DocURL, err))
			continue
		}
		text := ExtractFromHTML(html)
		requires := ExtractRequirements(text, target.Permission.APIKey)
		if len(requires) == 0 {
			continue
		}
		for i := range requires {
			requires[i].SourceURL = target.DocURL
		}
		key := target.Permission.APIKey + "\x00" + string(target.Access)
		edge, ok := edgeByKey[key]
		if !ok {
			edge = &model.DependencyEdge{
				Permission: target.Permission.APIKey,
				Access:     target.Access,
			}
			edgeByKey[key] = edge
		}
		edge.Requires = mergeRequires(edge.Requires, requires)
	}

	out := make([]model.DependencyEdge, 0, len(edgeByKey))
	for _, edge := range edgeByKey {
		sort.Slice(edge.Requires, func(i, j int) bool {
			if edge.Requires[i].Permission != edge.Requires[j].Permission {
				return edge.Requires[i].Permission < edge.Requires[j].Permission
			}
			return edge.Requires[i].Access < edge.Requires[j].Access
		})
		out = append(out, *edge)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Permission != out[j].Permission {
			return out[i].Permission < out[j].Permission
		}
		return out[i].Access < out[j].Access
	})
	return out, warnings, nil
}

func docURLsFromInventory(inv *inventory.Inventory) map[string]string {
	out := make(map[string]string)
	if inv == nil {
		return out
	}
	for perm, meta := range inv.Permissions {
		out[perm] = meta.DocURL
	}
	return out
}

func mergeRequires(existing, extra []model.RequiredPermission) []model.RequiredPermission {
	byKey := make(map[string]model.RequiredPermission, len(existing)+len(extra))
	for _, req := range existing {
		byKey[req.Permission+"\x00"+string(req.Access)] = req
	}
	for _, req := range extra {
		byKey[req.Permission+"\x00"+string(req.Access)] = req
	}
	out := make([]model.RequiredPermission, 0, len(byKey))
	for _, req := range byKey {
		out = append(out, req)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Permission != out[j].Permission {
			return out[i].Permission < out[j].Permission
		}
		return out[i].Access < out[j].Access
	})
	return out
}
