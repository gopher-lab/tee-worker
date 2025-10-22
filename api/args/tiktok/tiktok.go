package tiktok

import (
	"github.com/masa-finance/tee-worker/api/args/tiktok/query"
	"github.com/masa-finance/tee-worker/api/args/tiktok/transcription"
	"github.com/masa-finance/tee-worker/api/args/tiktok/trending"
)

type TranscriptionArguments = transcription.Arguments
type QueryArguments = query.Arguments
type TrendingArguments = trending.Arguments

var NewTranscriptionArguments = transcription.NewArguments
var NewQueryArguments = query.NewArguments
var NewTrendingArguments = trending.NewArguments
