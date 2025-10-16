package tiktok

import (
	"github.com/masa-finance/tee-worker/api/args/tiktok/query"
	"github.com/masa-finance/tee-worker/api/args/tiktok/transcription"
	"github.com/masa-finance/tee-worker/api/args/tiktok/trending"
)

type Transcription = transcription.Arguments
type Query = query.Arguments
type Trending = trending.Arguments
