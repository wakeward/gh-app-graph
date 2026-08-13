package scrape

import (
	"sort"

	"github.com/wakeward/github-app-permissions-graph/pkg/model"
)

// Target is one permission+access pair selected for prose scraping.
type Target struct {
	Permission model.Permission
	Access     model.AccessLevel
	Severity   model.Severity
	DocURL     string
}

// SelectTargets returns permission+access pairs to scrape. By default this is
// High/Critical severity only; widen includes every access level on rows
// flagged needs_investigation as well.
func SelectTargets(permissions []model.Permission, docURLs map[string]string, widen bool) []Target {
	var out []Target
	for _, p := range permissions {
		for _, al := range p.AccessLevels {
			if !includeAccess(p, al, widen) {
				continue
			}
			out = append(out, Target{
				Permission: p,
				Access:     al.Access,
				Severity:   al.Severity,
				DocURL:     docURLs[p.APIKey],
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Permission.APIKey != out[j].Permission.APIKey {
			return out[i].Permission.APIKey < out[j].Permission.APIKey
		}
		return out[i].Access < out[j].Access
	})
	return out
}

func includeAccess(p model.Permission, al model.AccessLevelDetail, widen bool) bool {
	if widen && p.NeedsInvestigation {
		return true
	}
	switch al.Severity {
	case model.SeverityCritical, model.SeverityHigh:
		return true
	default:
		return false
	}
}
