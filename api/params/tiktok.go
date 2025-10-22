package params

import (
	"github.com/masa-finance/tee-worker/api/args/tiktok"
)

type TikTokTranscription = Params[*tiktok.TranscriptionArguments]

type TikTokSearch = Params[*tiktok.QueryArguments]

type TikTokTrending = Params[*tiktok.TrendingArguments]

type TikTok = Generic
