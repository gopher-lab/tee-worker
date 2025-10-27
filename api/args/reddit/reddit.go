package reddit

import (
	"github.com/masa-finance/tee-worker/v2/api/args/reddit/search"
)

type SearchArguments = search.Arguments

var NewSearchArguments = search.NewArguments

var NewSearchPostsArguments = search.NewSearchPostsArguments
var NewSearchUsersArguments = search.NewSearchUsersArguments
var NewSearchCommunitiesArguments = search.NewSearchCommunitiesArguments
var NewScrapeUrlsArguments = search.NewScrapeUrlsArguments
