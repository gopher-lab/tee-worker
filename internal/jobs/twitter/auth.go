package twitter

import (
	"github.com/sirupsen/logrus"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	// Account-based auth
	Account *TwitterAccount
	BaseDir string
}

func NewScraper(config AuthConfig) *Scraper {

	// Fall back to account-based auth
	if config.Account == nil {
		logrus.Error("No authentication method provided")
		return nil
	}

	scraper := &Scraper{Scraper: newTwitterScraper()}

	// Try loading cookies
	if err := LoadCookies(scraper.Scraper, config.Account, config.BaseDir); err == nil {
		logrus.Debugf("Cookies loaded for user %s.", config.Account.Username)
		scraper.SetBearerToken()
		return scraper
	} else {
		logrus.WithError(err).Warnf("Failed to load cookies for user %s", config.Account.Username)
		return nil
	}
}
