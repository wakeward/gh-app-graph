package inventory

// octokitRaw mirrors the shape of octokit/app-permissions generated JSON.
type octokitRaw struct {
	Paths       map[string]map[string]octokitPathGrant `json:"paths"`
	Permissions map[string]octokitPermission           `json:"permissions"`
}

type octokitPathGrant struct {
	Access     string `json:"access"`
	Permission string `json:"permission"`
}

type octokitPermission struct {
	Read  []string `json:"read"`
	Write []string `json:"write"`
	URL   string   `json:"url"`
}
