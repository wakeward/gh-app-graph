package inventory

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/wakeward/gh-app-graph/pkg/fileio"
)

// FetchHTTP downloads raw octokit JSON from url.
func FetchHTTP(url string) ([]byte, error) {
	resp, err := http.Get(url) // #nosec G107 -- URL defaults to pinned octokit raw JSON; overridable for offline mirrors
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: HTTP %s", url, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return body, nil
}

// ReadFile loads raw octokit JSON from a local path (for offline/tests).
func ReadFile(path string) ([]byte, error) {
	body, err := os.ReadFile(path) // #nosec G304 -- path is a caller-supplied local inventory file
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return body, nil
}

// WriteJSON writes inv to path with stable, indented JSON encoding.
func WriteJSON(path string, inv *Inventory) error {
	data, err := json.MarshalIndent(inv, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal inventory: %w", err)
	}
	data = append(data, '\n')
	if err := fileio.Write(path, data); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// BuildFromRaw is the shared path: bytes in, normalized inventory out.
func BuildFromRaw(sourceURL string, raw []byte) (*Inventory, error) {
	return Normalize(sourceURL, raw, time.Now().UTC())
}
