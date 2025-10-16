package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/types"
)

var (
	ErrInvalidType       = errors.New("invalid type")
	ErrInvalidSort       = errors.New("invalid sort")
	ErrTimeInTheFuture   = errors.New("after field is in the future")
	ErrNoQueries         = errors.New("queries must be provided for all query types except scrapeurls")
	ErrNoUrls            = errors.New("urls must be provided for scrapeurls query type")
	ErrQueriesNotAllowed = errors.New("the scrapeurls query type does not admit queries")
	ErrUrlsNotAllowed    = errors.New("urls can only be provided for the scrapeurls query type")
	ErrUnmarshalling     = errors.New("failed to unmarshal reddit search arguments")
)

const (
	// These reflect the default values in https://apify.com/trudax/reddit-scraper/input-schema
	DefaultMaxItems       = 10
	DefaultMaxPosts       = 10
	DefaultMaxComments    = 10
	DefaultMaxCommunities = 2
	DefaultMaxUsers       = 2
	DefaultSort           = types.RedditSortNew
)

const DomainSuffix = "reddit.com"

// Verify interface implementation
var _ base.JobArgument = (*Arguments)(nil)

// Arguments defines args for Reddit scrapes
// see https://apify.com/trudax/reddit-scraper
type Arguments struct {
	Type           types.Capability     `json:"type"`
	Queries        []string             `json:"queries"`
	URLs           []string             `json:"urls"`
	Sort           types.RedditSortType `json:"sort"`
	IncludeNSFW    bool                 `json:"include_nsfw"`
	SkipPosts      bool                 `json:"skip_posts"`      // Valid only for searchusers
	After          time.Time            `json:"after"`           // valid only for scrapeurls and searchposts
	MaxItems       uint                 `json:"max_items"`       // Max number of items to scrape (total), default 10
	MaxResults     uint                 `json:"max_results"`     // Max number of results per page, default MaxItems
	MaxPosts       uint                 `json:"max_posts"`       // Max number of posts per page, default 10
	MaxComments    uint                 `json:"max_comments"`    // Max number of comments per page, default 10
	MaxCommunities uint                 `json:"max_communities"` // Max number of communities per page, default 2
	MaxUsers       uint                 `json:"max_users"`       // Max number of users per page, default 2
	NextCursor     string               `json:"next_cursor"`
}

func (t *Arguments) UnmarshalJSON(data []byte) error {
	type Alias Arguments
	aux := &struct{ *Alias }{Alias: (*Alias)(t)}
	if err := json.Unmarshal(data, aux); err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshalling, err)
	}
	t.SetDefaultValues()
	return t.Validate()
}

// SetDefaultValues sets the default values for the parameters that were not provided and canonicalizes the strings for later validation
func (r *Arguments) SetDefaultValues() {
	if r.MaxItems == 0 {
		r.MaxItems = DefaultMaxItems
	}
	if r.MaxPosts == 0 {
		r.MaxPosts = DefaultMaxPosts
	}
	if r.MaxComments == 0 {
		r.MaxComments = DefaultMaxComments
	}
	if r.MaxCommunities == 0 {
		r.MaxCommunities = DefaultMaxCommunities
	}
	if r.MaxUsers == 0 {
		r.MaxUsers = DefaultMaxUsers
	}
	if r.MaxItems != 0 {
		r.MaxResults = r.MaxItems
	} else if r.MaxResults == 0 {
		r.MaxResults = DefaultMaxItems
	}
	if r.Sort == "" {
		r.Sort = DefaultSort
	}

	r.Sort = types.RedditSortType(strings.ToLower(string(r.Sort)))
}

func (r *Arguments) Validate() error {
	var errs []error

	if !types.AllRedditQueryTypes.Contains(r.Type) {
		errs = append(errs, ErrInvalidType)
	}

	if !types.AllRedditSortTypes.Contains(r.Sort) {
		errs = append(errs, ErrInvalidSort)
	}

	if time.Now().Before(r.After) {
		errs = append(errs, ErrTimeInTheFuture)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	if r.Type == types.CapScrapeUrls {
		if len(r.URLs) == 0 {
			errs = append(errs, ErrNoUrls)
		}
		if len(r.Queries) > 0 {
			errs = append(errs, ErrQueriesNotAllowed)
		}

		for _, u := range r.URLs {
			u, err := url.Parse(u)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s is not a valid URL", u))
			} else {
				if !strings.HasSuffix(strings.ToLower(u.Host), DomainSuffix) {
					errs = append(errs, fmt.Errorf("invalid Reddit URL %s", u))
				}
				if !strings.HasPrefix(u.Path, "/r/") {
					errs = append(errs, fmt.Errorf("%s is not a Reddit post or comment URL (missing /r/)", u))
				}
				if !strings.Contains(u.Path, "/comments/") {
					errs = append(errs, fmt.Errorf("%s is not a Reddit post or comment URL (missing /comments/)", u))
				}
			}
		}
	} else {
		if len(r.Queries) == 0 {
			errs = append(errs, ErrNoQueries)
		}
		if len(r.URLs) > 0 {
			errs = append(errs, ErrUrlsNotAllowed)
		}
	}

	return errors.Join(errs...)
}

// GetCapability returns the capability of the arguments
func (r *Arguments) GetCapability() types.Capability {
	return r.Type
}

// ValidateCapability validates the capability of the arguments
func (r *Arguments) ValidateCapability(jobType types.JobType) error {
	return jobType.ValidateCapability(&r.Type)
}

// NewArguments creates a new Arguments instance with the specified capability
// and applies default values immediately
func NewArguments(capability types.Capability) Arguments {
	args := Arguments{
		Type: capability,
	}
	args.SetDefaultValues()
	return args
}

// NewSearchPostsArguments creates a new Arguments instance for searching posts
func NewSearchPostsArguments() Arguments {
	return NewArguments(types.CapSearchPosts)
}

// NewSearchUsersArguments creates a new Arguments instance for searching users
func NewSearchUsersArguments() Arguments {
	return NewArguments(types.CapSearchUsers)
}

// NewSearchCommunitiesArguments creates a new Arguments instance for searching communities
func NewSearchCommunitiesArguments() Arguments {
	return NewArguments(types.CapSearchCommunities)
}

// NewScrapeUrlsArguments creates a new Arguments instance for scraping URLs
func NewScrapeUrlsArguments() Arguments {
	return NewArguments(types.CapScrapeUrls)
}
