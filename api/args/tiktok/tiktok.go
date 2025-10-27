package tiktok

import (
	"github.com/masa-finance/tee-worker/v2/api/args/tiktok/query"
	"github.com/masa-finance/tee-worker/v2/api/args/tiktok/transcription"
	"github.com/masa-finance/tee-worker/v2/api/args/tiktok/trending"
)

type TranscriptionArguments = transcription.Arguments
type QueryArguments = query.Arguments
type TrendingArguments = trending.Arguments

var NewTranscriptionArguments = transcription.NewArguments
var NewQueryArguments = query.NewArguments
var NewTrendingArguments = trending.NewArguments
