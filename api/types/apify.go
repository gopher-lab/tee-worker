package types

type ApifyProxy struct {
	UseApifyProxy    bool     `json:"useApifyProxy"`
	ApifyProxyGroups []string `json:"apifyProxyGroups"`
}
