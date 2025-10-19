package capabilities

import (
	"slices"
	"strings"

	"maps"

	"github.com/masa-finance/tee-worker/api/types"
	"github.com/masa-finance/tee-worker/internal/apify"
	"github.com/masa-finance/tee-worker/internal/config"
	"github.com/masa-finance/tee-worker/internal/jobs/twitter"
	"github.com/masa-finance/tee-worker/pkg/client"
	util "github.com/masa-finance/tee-worker/pkg/util"
	"github.com/sirupsen/logrus"
)

// JobServerInterface defines the methods we need from JobServer to avoid circular dependencies
type JobServerInterface interface {
	GetWorkerCapabilities() types.WorkerCapabilities
}

// DetectCapabilities automatically detects available capabilities based on configuration
// Always performs real capability detection by probing APIs and actors to ensure accurate reporting
func DetectCapabilities(jc config.JobConfiguration, jobServer JobServerInterface) types.WorkerCapabilities {
	// Always perform real capability detection to ensure accurate reporting
	// This guarantees miners report only capabilities they actually have access to
	capabilities := make(types.WorkerCapabilities)

	// Start with always available capabilities
	maps.Copy(capabilities, types.AlwaysAvailableCapabilities)

	// Check what Twitter authentication methods are available
	accounts := jc.GetStringSlice("twitter_accounts", nil)
	apiKeys := jc.GetStringSlice("twitter_api_keys", nil)
	apifyApiKey := jc.GetString("apify_api_key", "")
	geminiApiKey := config.LlmApiKey(jc.GetString("gemini_api_key", ""))
	claudeApiKey := config.LlmApiKey(jc.GetString("claude_api_key", ""))

	hasAccounts := len(accounts) > 0
	hasApiKeys := len(apiKeys) > 0
	hasApifyKey := hasValidApifyKey(apifyApiKey)
	hasLLMKey := geminiApiKey.IsValid() || claudeApiKey.IsValid()

	// Add Twitter capabilities based on available authentication
	var twitterCaps []types.Capability

	// Add credential-based capabilities if we have accounts
	if hasAccounts {
		twitterCaps = append(twitterCaps,
			types.CapSearchByQuery,
			types.CapSearchByProfile,
			types.CapGetById,
			types.CapGetReplies,
			types.CapGetRetweeters,
			types.CapGetMedia,
			types.CapGetProfileById,
			types.CapGetTrends,
			types.CapGetSpace,
			types.CapGetProfile,
			types.CapGetTweets,
		)
	}

	// Add API-based capabilities if we have API keys
	if hasApiKeys {
		// Check for elevated API capabilities
		if hasElevatedApiKey(apiKeys) {
			twitterCaps = append(twitterCaps, types.CapSearchByFullArchive)
		}
	}

	// Only add capabilities if we have any supported capabilities
	if len(twitterCaps) > 0 {
		capabilities[types.TwitterJob] = twitterCaps
	}

	if hasApifyKey {
		// Create an Apify client for probing actors
		c, err := client.NewApifyClient(apifyApiKey)
		if err != nil {
			logrus.Errorf("Failed to create Apify client for access probes: %v", err)
		} else {
			// Aggregate capabilities per job from accessible actors
			jobToSet := map[types.JobType]*util.Set[types.Capability]{}

			for _, actor := range apify.Actors {
				// Web requires a valid Gemini API key
				if actor.JobType == types.WebJob && !hasLLMKey {
					logrus.Debug("Skipping Web actor due to missing Gemini key")
					continue
				}

				if ok, _ := c.ProbeActorAccess(actor.ActorId, actor.DefaultInput); ok {
					if _, exists := jobToSet[actor.JobType]; !exists {
						jobToSet[actor.JobType] = util.NewSet[types.Capability]()
					}
					jobToSet[actor.JobType].Add(actor.Capabilities...)
				} else {
					logrus.Warnf("Apify token does not have access to actor %s", actor.ActorId)
				}
			}

			// Union accessible-actor caps into existing caps
			for job, set := range jobToSet {
				existingCaps := util.NewSet(capabilities[job]...)
				capabilities[job] = existingCaps.Add(set.Items()...).Items()
			}
		}
	}

	return capabilities
}

// hasElevatedApiKey checks if any of the provided API keys are elevated
func hasElevatedApiKey(apiKeys []string) bool {
	if len(apiKeys) == 0 {
		return false
	}

	// Parse API keys and create account manager to detect types
	parsedApiKeys := parseApiKeys(apiKeys)
	accountManager := twitter.NewTwitterAccountManager(nil, parsedApiKeys)

	// Detect all API key types
	accountManager.DetectAllApiKeyTypes()

	// Check if any key is elevated
	return slices.ContainsFunc(accountManager.GetApiKeys(), func(apiKey *twitter.TwitterApiKey) bool {
		return apiKey.Type == twitter.TwitterApiKeyTypeElevated
	})
}

// parseApiKeys converts string API keys to TwitterApiKey structs
func parseApiKeys(apiKeys []string) []*twitter.TwitterApiKey {
	result := make([]*twitter.TwitterApiKey, 0, len(apiKeys))
	for _, key := range apiKeys {
		if trimmed := strings.TrimSpace(key); trimmed != "" {
			result = append(result, &twitter.TwitterApiKey{
				Key: trimmed,
			})
		}
	}
	return result
}

// hasValidApifyKey checks if the provided Apify API key is valid by attempting to validate it
func hasValidApifyKey(apifyApiKey string) bool {
	if apifyApiKey == "" {
		return false
	}

	// Create temporary Apify client and validate the key
	apifyClient, err := client.NewApifyClient(apifyApiKey)
	if err != nil {
		logrus.Errorf("Failed to create Apify client during capability detection: %v", err)
		return false
	}

	if err := apifyClient.ValidateApiKey(); err != nil {
		logrus.Errorf("Apify API key validation failed during capability detection: %v", err)
		return false
	}

	logrus.Infof("Apify API key validated successfully during capability detection")
	return true
}
