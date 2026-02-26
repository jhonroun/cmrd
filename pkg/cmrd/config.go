package cmrd

import (
	"strings"
	"time"
)

// Config configures library behavior.
type Config struct {
	Aria2Path            string
	DownloadDir          string
	Proxy                string
	ProxyAuth            string
	HTTPTimeout          time.Duration
	DeleteInputAfterDone bool
}

// DefaultConfig returns recommended defaults.
func DefaultConfig() Config {
	return Config{
		DownloadDir:          "downloads",
		HTTPTimeout:          30 * time.Second,
		DeleteInputAfterDone: true,
	}
}

func (c Config) normalized() Config {
	cfg := c
	if strings.TrimSpace(cfg.DownloadDir) == "" {
		cfg.DownloadDir = "downloads"
	}
	if cfg.HTTPTimeout <= 0 {
		cfg.HTTPTimeout = 30 * time.Second
	}
	return cfg
}
