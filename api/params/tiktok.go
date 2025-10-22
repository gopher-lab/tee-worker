package params

import (
	"github.com/masa-finance/tee-worker/api/args/tiktok"
	"github.com/masa-finance/tee-worker/api/types"
)

type TikTokTranscription = Params[*tiktok.TranscriptionArguments]

type TikTokSearch = Params[*tiktok.QueryArguments]

type TikTokTrending = Params[*tiktok.TrendingArguments]

type TikTokGenericArgs struct {
	Data map[string]any `json:",inline"`
}

type TikTok = Params[TikTokGenericArgs]

func (t TikTokGenericArgs) GetCapability() types.Capability {
	return types.CapEmpty
}

func (t TikTokGenericArgs) SetDefaultValues() {
}

func (t TikTokGenericArgs) Validate() error {
	return nil
}
