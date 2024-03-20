package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/lmittmann/tint"

	"go.xsfx.dev/glucose_exporter/httpslog"
	"go.xsfx.dev/glucose_exporter/internal/cache"
	"go.xsfx.dev/glucose_exporter/internal/config"
	"go.xsfx.dev/glucose_exporter/internal/metrics"
)

const addr = ":2112"

func main() {
	logOpts := &tint.Options{Level: slog.LevelInfo, TimeFormat: time.Kitchen}
	initLogging(logOpts)

	if err := env.Parse(&config.Cfg); err != nil {
		slog.Error("parsing env config", "err", err)
		os.Exit(1)
	}

	if config.Cfg.Debug {
		initLogging(&tint.Options{Level: slog.LevelDebug, TimeFormat: time.Kitchen})
	}

	cacheFileLogger := slog.With("file", cache.FullPath())

	// Check if cache files needs to be created.
	_, err := os.Stat(cache.FullPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(cache.FullPath(), []byte("{}"), 0o600); err != nil {
				cacheFileLogger.Error("init cache file", "err", err)
				os.Exit(1)
			} else {
				cacheFileLogger.Debug("init cache file")
			}
		} else {
			cacheFileLogger.Error("getting cache file stat", "err", err)
			os.Exit(1)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", metrics.Handler)

	slog.Info("listening", "addr", addr)

	if err := http.ListenAndServe(addr, httpslog.Handler()(mux)); err != nil {
		slog.Error("listen and serve", "err", err)
		os.Exit(1)
	}
}

func initLogging(opts *tint.Options) {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, opts)))
}
