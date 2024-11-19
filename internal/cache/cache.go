package cache

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path"

	"dario.cat/mergo"
	"go.xsfx.dev/glucose_exporter/internal/config"
	"go.xsfx.dev/glucose_exporter/internal/epoch"
)

const cacheFile = "cache.json"

type Cache struct {
	JWT       string      `json:"jwt,omitempty"`
	Expires   epoch.Epoch `json:"expires,omitempty"`
	BaseURL   string      `json:"base_url,omitempty"`
	AccountID string      `json:"account_id,omitempty"`
}

func FullPath() string {
	return path.Join(config.Cfg.CacheDir, cacheFile)
}

func Load() (Cache, error) {
	var c Cache

	slog.Debug("reading cache", "file", FullPath())

	b, err := os.ReadFile(FullPath())
	if err != nil {
		return Cache{}, fmt.Errorf("reading cache file: %w", err)
	}

	if err := json.Unmarshal(b, &c); err != nil {
		return Cache{}, fmt.Errorf("unmarshal cache: %w", err)
	}

	return c, nil
}

func Save(c Cache) error {
	slog.Debug("writing cache", "file", FullPath())

	s, err := Load()
	if err != nil {
		return fmt.Errorf("loading cache: %w", err)
	}

	if err := mergo.Merge(&s, c, mergo.WithOverride); err != nil {
		return fmt.Errorf("merging cache: %w", err)
	}

	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling cache: %w", err)
	}

	if err := os.WriteFile(FullPath(), b, 0o600); err != nil {
		return fmt.Errorf("writing cache file: %w", err)
	}

	return nil
}
