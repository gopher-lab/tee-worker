package jobs

import (
	"github.com/masa-finance/tee-worker/api/args/tiktok"
)

type TikTokTranscriptionParams = Params[*tiktok.TranscriptionArguments]

type TikTokSearchParams = Params[*tiktok.QueryArguments]

type TikTokTrendingParams = Params[*tiktok.TrendingArguments]
