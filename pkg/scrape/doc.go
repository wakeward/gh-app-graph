// Package scrape implements the Type B (prose-scraped) dependency miner:
// for permissions flagged High/Critical severity, it scrapes each
// endpoint's own GitHub docs page (via goquery) for free-text statements of
// an additional required permission, and writes low-confidence,
// needs_review draft edges to
// data/generated/scraped-dependencies-draft.yaml.
//
// TODO(scrape-prose): implement the scraper, scoped to High/Critical
// severity permissions only per the plan.
package scrape
