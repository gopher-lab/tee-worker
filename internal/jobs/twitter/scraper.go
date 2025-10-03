package twitter

import (
	twitterscraper "github.com/imperatrona/twitter-scraper"
)

type Scraper struct {
	*twitterscraper.Scraper
	skipLoginVerification bool // false by default
}

func newTwitterScraper() *twitterscraper.Scraper {
	return twitterscraper.New()
}

// SetBearerToken sets the bearer token via scraper's IsLoggedIn check
func (s *Scraper) SetBearerToken() bool {
	s.Scraper.IsLoggedIn()
	return true
}
