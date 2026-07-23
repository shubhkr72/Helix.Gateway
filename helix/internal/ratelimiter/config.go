package ratelimiter

type Config struct {
	Capacity    float64 `yaml:"capacity"`
	RefillRate  float64 `yaml:"refill_rate"`
	KeyStrategy string  `yaml:"key_strategy"`
}
