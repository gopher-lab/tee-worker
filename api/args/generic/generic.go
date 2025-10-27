package generic

import (
	"encoding/json"
	"maps"

	"github.com/masa-finance/tee-worker/v2/api/args/base"
	"github.com/masa-finance/tee-worker/v2/api/types"
)

type Arguments struct {
	base.Arguments
	Data map[string]any
}

func (g *Arguments) UnmarshalJSON(data []byte) error {
	// First unmarshal into a map to get all fields
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract the type field
	if typeVal, ok := raw["type"]; ok {
		if typeStr, ok := typeVal.(string); ok {
			g.Type = types.Capability(typeStr)
		}
		delete(raw, "type") // Remove type from raw map
	}

	// Store all other fields in Data
	g.Data = raw
	return nil
}

func (g Arguments) MarshalJSON() ([]byte, error) {
	// Combine Type and Data into a single map
	result := make(map[string]any)
	maps.Copy(result, g.Data)
	result["type"] = g.Type
	return json.Marshal(result)
}
