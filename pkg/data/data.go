package data

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/wakeward/gh-app-graph/pkg/model"
)

//go:embed bundled/toxic-combinations.json
var toxicCombinationsJSON []byte

//go:embed bundled/permissions.resolved.json
var permissionsResolvedJSON []byte

// LoadToxicCombinations returns compiled-in toxic combination rules from the
// last cmd/build output copied into pkg/data/bundled/.
func LoadToxicCombinations() ([]model.ToxicCombination, error) {
	var combos []model.ToxicCombination
	if err := json.Unmarshal(toxicCombinationsJSON, &combos); err != nil {
		return nil, fmt.Errorf("decode embedded toxic-combinations.json: %w", err)
	}
	return combos, nil
}

// LoadResolvedPermissions returns the compiled-in resolved permission catalog.
func LoadResolvedPermissions() ([]model.Permission, error) {
	var envelope struct {
		Permissions []model.Permission `json:"permissions"`
	}
	if err := json.Unmarshal(permissionsResolvedJSON, &envelope); err != nil {
		return nil, fmt.Errorf("decode embedded permissions.resolved.json: %w", err)
	}
	return envelope.Permissions, nil
}
