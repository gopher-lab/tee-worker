package params

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/masa-finance/tee-worker/api/types"
)

type HybridQuery struct {
	Query  string  `json:"query"`
	Weight float64 `json:"weight"`
}

var _ JobParameters = (*HybridSearch)(nil)

// HybridSearch defines parameters for hybrid search
// TODO: At some point we could replace `TextQuery` and `SimilarityQuery` with a single slice that receives N queries. The issue right now is that, because of the Milvus API, we can't have an arbitrary number of queries (see https://github.com/milvus-io/milvus/issues/41261). Once that issue is resolved we can fix this.
type HybridSearch struct {
	TextQuery       HybridQuery    `json:"text_query"`       // Optional, either TextQuery or SimilarityQuery must be specified
	SimilarityQuery HybridQuery    `json:"similarity_query"` // Mandatory, 1 or more queries to execute
	Keywords        []string       `json:"keywords"`         // Optional, keywords to filter for in keyword search
	Operator        string         `json:"keyword_operator"` // Optional, operator ("and" / "or") to use in keyword search. Default is "and"
	MaxResults      int            `json:"max_results"`      // Optional, max number of results
	Sources         []types.Source `json:"sources"`
}

// Validate validates the hybrid search parameters
func (t HybridSearch) Validate(cfg *types.SearchConfig) error {
	if t.TextQuery.Weight <= 0 || t.TextQuery.Weight > 1 || t.SimilarityQuery.Weight <= 0 || t.SimilarityQuery.Weight > 1 {
		return fmt.Errorf("weights must be greater than or equal to 0, and less than 1, got %f and %f", t.TextQuery.Weight, t.SimilarityQuery.Weight)
	}

	op := strings.ToLower(t.Operator)
	if op != "and" && op != "or" && op != "" {
		return fmt.Errorf(`keyword_operator must be "and", "or" or "", not "%s"`, t.Operator)
	}

	for _, s := range t.Sources {
		if !slices.Contains(types.Sources, s) {
			return fmt.Errorf("source must be one of %v, got %s", types.Sources, s)
		}
	}

	return nil
}

func (t HybridSearch) Timeout() time.Duration {
	return 0
}

func (t HybridSearch) PollInterval() time.Duration {
	return 0
}

func (t HybridSearch) Arguments(cfg *types.SearchConfig) map[string]any {
	t.ApplyDefaults(cfg)

	return map[string]any{
		"text_query":       t.TextQuery,
		"similarity_query": t.SimilarityQuery,
		"keywords":         t.Keywords,
		"operator":         t.Operator,
		"max_results":      t.MaxResults,
		"sources":          t.Sources,
	}
}

func (t HybridSearch) Type() types.JobType {
	return "hybrid-search"
}

// ApplyDefaults applies default values to the hybrid search parameters
func (t *HybridSearch) ApplyDefaults(cfg *types.SearchConfig) {
	switch {
	case t.MaxResults == 0:
		t.MaxResults = int(cfg.DefaultMaxResults)
	case t.MaxResults < int(cfg.MinMaxResults):
		t.MaxResults = int(cfg.MinMaxResults)
	case t.MaxResults > int(cfg.MaxMaxResults):
		t.MaxResults = int(cfg.MaxMaxResults)
	}
}
