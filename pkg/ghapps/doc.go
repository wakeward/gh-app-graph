// Package ghapps is a thin client for the public, unauthenticated
// GET /apps/{app_slug} endpoint, used by cmd/evaluate-app to fetch a
// third-party app's default requested permissions without needing org
// access or a token. An optional token may still be supplied for a higher
// rate limit.
//
// TODO(evaluate-app): implement the client and response mapping into
// model.AppPermissionSet.
package ghapps
