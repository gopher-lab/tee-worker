package params

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/masa-finance/tee-worker/api/types"
)

var _ JobParameters = (*SimilaritySearchParams)(nil)

type SimilaritySearchParams struct {
	Query           string         `json:"query"`            // Mandatory, query for similarity search in keyword search
	Keywords        []string       `json:"keywords"`         // Optional, keywords to filter for in keyword search
	KeywordOperator string         `json:"keyword_operator"` // Optional, operator ("and" / "or") to use in keyword search. Default is "and"
	Sources         []types.Source `json:"sources"`          // Optional, sources to query
	MaxResults      int            `json:"max_results"`      // Optional, max number of results for keyword search
}

func (t SimilaritySearchParams) Validate(cfg *SearchConfig) error {
	if t.Query == "" {
		return errors.New("query is required")
	}

	t.KeywordOperator = strings.ToLower(t.KeywordOperator)
	if t.KeywordOperator != "and" && t.KeywordOperator != "or" && t.KeywordOperator != "" {
		return fmt.Errorf(`keyword_operator must be "and", "or" or "", not "%s"`, t.KeywordOperator)
	}

	for _, s := range t.Sources {
		if !slices.Contains(types.Sources, s) {
			return fmt.Errorf("source must be one of %v, got %s", types.Sources, s)
		}
	}

	return nil
}

func (t SimilaritySearchParams) Timeout() time.Duration {
	return 0
}

func (t SimilaritySearchParams) PollInterval() time.Duration {
	return 0
}

func (t SimilaritySearchParams) Type() types.JobType {
	return "similarity-search"
}

func (t SimilaritySearchParams) Arguments(cfg *SearchConfig) map[string]any {
	t.ApplyDefaults(cfg)

	return map[string]any{
		"query":            t.Query,
		"keywords":         t.Keywords,
		"keyword_operator": strings.ToLower(t.KeywordOperator),
		"max_results":      t.MaxResults,
		"sources":          t.Sources,
	}
}

func (t *SimilaritySearchParams) ApplyDefaults(cfg *SearchConfig) {
	switch {
	case t.MaxResults == 0:
		t.MaxResults = int(cfg.DefaultMaxResults)
	case t.MaxResults < int(cfg.MinMaxResults):
		t.MaxResults = int(cfg.MinMaxResults)
	case t.MaxResults > int(cfg.MaxMaxResults):
		t.MaxResults = int(cfg.MaxMaxResults)
	}

	if t.KeywordOperator == "" {
		t.KeywordOperator = "and"
	}
	t.KeywordOperator = strings.ToLower(t.KeywordOperator)
}
