package scrape

import (
	"regexp"
	"strings"

	"github.com/wakeward/github-app-permissions-graph/pkg/model"
)

var (
	grantPattern = regexp.MustCompile("(?i)`([a-z][a-z0-9_]*):(?:read|write)`")
	namePattern  = regexp.MustCompile("(?i)`([a-z][a-z0-9_]*)` permission")
)

// ExtractRequirements finds permission names referenced in GitHub docs prose.
// sourcePermission is excluded from results.
func ExtractRequirements(text, sourcePermission string) []model.RequiredPermission {
	text = strings.ToLower(text)
	seen := make(map[string]struct{})
	var out []model.RequiredPermission

	add := func(permission string, access model.AccessLevel, condition string) {
		if permission == "" || permission == sourcePermission {
			return
		}
		key := permission + "\x00" + string(access)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		out = append(out, model.RequiredPermission{
			Permission:  permission,
			Access:      access,
			Condition:   condition,
			Confidence:  model.ConfidenceLow,
			Type:        model.DependencyTypeProseScraped,
			NeedsReview: true,
		})
	}

	for _, m := range grantPattern.FindAllStringSubmatch(text, -1) {
		access := model.AccessRead
		if strings.Contains(m[0], ":write") {
			access = model.AccessWrite
		}
		add(m[1], access, "Prose mentions `"+m[1]+":"+string(access)+"`")
	}
	for _, m := range namePattern.FindAllStringSubmatch(text, -1) {
		add(m[1], model.AccessWrite, "Prose mentions `"+m[1]+"` permission")
	}
	return out
}

// ExtractFromHTML pulls visible text from a GitHub docs HTML page.
func ExtractFromHTML(html string) string {
	// Lightweight extraction keeps tests fast without pulling in goquery for
	// every package consumer. strip tags and collapse whitespace.
	re := regexp.MustCompile(`(?s)<(script|style)[^>]*>.*?</(script|style)>`)
	text := re.ReplaceAllString(html, " ")
	text = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}
