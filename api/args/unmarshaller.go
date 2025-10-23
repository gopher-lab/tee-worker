package args

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/masa-finance/tee-worker/api/args/base"
	"github.com/masa-finance/tee-worker/api/args/linkedin"
	"github.com/masa-finance/tee-worker/api/args/reddit"
	"github.com/masa-finance/tee-worker/api/args/telemetry"
	"github.com/masa-finance/tee-worker/api/args/tiktok"
	"github.com/masa-finance/tee-worker/api/args/twitter"
	"github.com/masa-finance/tee-worker/api/args/web"
	"github.com/masa-finance/tee-worker/api/types"
)

var (
	ErrUnknownJobType    = errors.New("unknown job type")
	ErrUnknownCapability = errors.New("unknown capability")
	ErrFailedToUnmarshal = errors.New("failed to unmarshal job arguments")
	ErrFailedToMarshal   = errors.New("failed to marshal job arguments")
)

type Args = map[string]any

// UnmarshalJobArguments unmarshals job arguments from a generic map into the appropriate typed struct
// This works with both tee-indexer and tee-worker JobArgument types
func UnmarshalJobArguments(jobType types.JobType, args Args) (base.JobArgument, error) {
	switch jobType {
	case types.WebJob:
		return unmarshalWebArguments(args)

	case types.TiktokJob:
		return unmarshalTikTokArguments(args)

	case types.TwitterJob:
		return unmarshalTwitterArguments(args)

	case types.LinkedInJob:
		return unmarshalLinkedInArguments(args)

	case types.RedditJob:
		return unmarshalRedditArguments(args)

	case types.TelemetryJob:
		return unmarshalTelemetryArguments(args)

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownJobType, jobType)
	}
}

// Helper functions for unmarshaling specific argument types
func unmarshalWebArguments(args Args) (*web.ScraperArguments, error) {
	webArgs := &web.ScraperArguments{}
	if err := unmarshalToStruct(args, webArgs); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return webArgs, nil
}

func unmarshalTikTokArguments(args Args) (base.JobArgument, error) {
	minimal := base.Arguments{}
	if err := unmarshalToStruct(args, &minimal); err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	switch minimal.Type {
	case types.CapSearchByQuery:
		searchArgs := &tiktok.QueryArguments{}
		if err := unmarshalToStruct(args, searchArgs); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		return searchArgs, nil
	case types.CapSearchByTrending:
		searchArgs := &tiktok.TrendingArguments{}
		if err := unmarshalToStruct(args, searchArgs); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		return searchArgs, nil
	case types.CapTranscription:
		transcriptionArgs := &tiktok.TranscriptionArguments{}
		if err := unmarshalToStruct(args, transcriptionArgs); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		return transcriptionArgs, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownCapability, minimal.Type)
	}
}

func unmarshalTwitterArguments(args Args) (*twitter.SearchArguments, error) {
	twitterArgs := &twitter.SearchArguments{}
	if err := unmarshalToStruct(args, twitterArgs); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToUnmarshal, err)
	}
	return twitterArgs, nil
}

func unmarshalLinkedInArguments(args Args) (*linkedin.ProfileArguments, error) {
	linkedInArgs := &linkedin.ProfileArguments{}
	if err := unmarshalToStruct(args, linkedInArgs); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToUnmarshal, err)
	}
	return linkedInArgs, nil
}

func unmarshalRedditArguments(args Args) (*reddit.SearchArguments, error) {
	redditArgs := &reddit.SearchArguments{}
	if err := unmarshalToStruct(args, redditArgs); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToUnmarshal, err)
	}
	return redditArgs, nil
}

func unmarshalTelemetryArguments(args Args) (*telemetry.Arguments, error) {
	telemetryArgs := &telemetry.Arguments{}
	if err := unmarshalToStruct(args, telemetryArgs); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToUnmarshal, err)
	}
	return telemetryArgs, nil
}

// unmarshalToStruct converts a map[string]any to a struct using JSON marshal/unmarshal
func unmarshalToStruct(args Args, target any) error {
	data, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToMarshal, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUnmarshal, err)
	}

	return nil
}
