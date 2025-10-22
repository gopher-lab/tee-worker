package types

type SearchConfig struct {
	DefaultMaxResults uint `env:"DEFAULT_MAX_RESULTS,default=20"`
	MinMaxResults     uint `env:"MIN_MAX_RESULTS,default=1"`
	MaxMaxResults     uint `env:"MAX_MAX_RESULTS,default=100"`
}
