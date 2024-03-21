//nolint:forbidigo,gochecknoglobals
package main

import (
	"errors"
	"flag"
	"fmt"
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

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const help = "expected 'serve' or 'version' subcommands"

func main() {
	logOpts := &tint.Options{Level: slog.LevelInfo, TimeFormat: time.Kitchen}
	initLogging(logOpts)

	if len(os.Args) < 2 {
		slog.Error(help)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		flags := flag.NewFlagSet("serve", flag.ExitOnError)

		addr := flags.String("addr", ":2112", "address to listen")
		flags.BoolFunc("help", "print help", func(_ string) error {
			flags.PrintDefaults()
			os.Exit(0)

			return nil
		})

		if err := flags.Parse(os.Args[2:]); err != nil {
			slog.Error("parsing flags", "err", err)
			os.Exit(0)
		}

		serveCmd(*addr)

	case "version":
		versionCmd()

	default:
		slog.Error(help)
		os.Exit(1)
	}
}

func serveCmd(addr string) {
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
			}

			cacheFileLogger.Debug("init cache file")
		} else {
			cacheFileLogger.Error("getting cache file stat", "err", err)
			os.Exit(1)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", metrics.Handler)

	slog.Info("listening", "addr", addr)

	server := &http.Server{
		Addr:              addr,
		ReadTimeout:       time.Second,
		ReadHeaderTimeout: time.Second,
		Handler:           httpslog.Handler()(mux),
	}

	if err := server.ListenAndServe(); err != nil {
		slog.Error("listen and serve", "err", err)
		os.Exit(1)
	}
}

func versionCmd() {
	fmt.Printf("glucose_exporter %s, commit %s, %s\n", version, commit, date)
}

func initLogging(opts *tint.Options) {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, opts)))
}
