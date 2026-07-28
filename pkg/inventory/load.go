package inventory

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadJSON reads a normalized Inventory from path.
func LoadJSON(path string) (*Inventory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var inv Inventory
	if err := json.Unmarshal(data, &inv); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	if len(inv.EndpointIndex) == 0 {
		return nil, fmt.Errorf("%s: empty endpoint_index", path)
	}
	return &inv, nil
}
