// Package data go:embeds the built, resolved JSON artifacts
// (permissions.resolved.json, toxic-combinations.json) directly into any
// binary that imports this module, so a consumer like gh-app-check can call
// data.LoadToxicCombinations() and get compiled-in data with no runtime
// file I/O or network fetch, and no dependency on this repo's checkout
// layout.
//
// TODO(build-render): once build.py's Go equivalent (cmd/build) produces
// permissions.resolved.json and toxic-combinations.json, add the
// //go:embed directives and Load* functions here.
package data
