package cache_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/matryer/is"
	"go.xsfx.dev/glucose_exporter/internal/cache"
	"go.xsfx.dev/glucose_exporter/internal/config"
	"go.xsfx.dev/glucose_exporter/internal/epoch"
)

func TestLoad(t *testing.T) {
	is := is.New(t)

	d := t.TempDir()

	config.Cfg.CacheDir = d

	t.Cleanup(func() { config.Cfg.CacheDir = "" })

	err := os.WriteFile(
		path.Join(d, "cache.json"),
		[]byte(`
			{
				"jwt": "f00b4r",
				"expires": 1726302083
			}
		`),
		0o600,
	)
	is.NoErr(err)

	c, err := cache.Load()
	is.NoErr(err)

	is.Equal(c.JWT, "f00b4r")
	is.True(time.Time(c.Expires).Equal(time.Date(2024, 9, 14, 8, 21, 23, 0, time.UTC)))
}

func TestSave(t *testing.T) {
	tables := []struct {
		name     string
		prefill  []byte
		input    cache.Cache
		expected cache.Cache
	}{
		{
			name:     "empty init",
			prefill:  []byte(`{}`),
			input:    cache.Cache{JWT: "f00b4r"},
			expected: cache.Cache{JWT: "f00b4r"},
		},
		{
			name:    "adding jwt",
			prefill: []byte(`{"jwt":"b4rf00"}`),
			input: cache.Cache{
				Expires: epoch.Epoch(time.Date(2024, 9, 14, 8, 21, 23, 0, time.UTC)),
			},
			expected: cache.Cache{
				JWT:     "b4rf00",
				Expires: epoch.Epoch(time.Date(2024, 9, 14, 8, 21, 23, 0, time.UTC)),
			},
		},
		{
			name:    "changing jwt",
			prefill: []byte(`{"jwt":"b4rf00"}`),
			input: cache.Cache{
				JWT: "f00b4r",
			},
			expected: cache.Cache{
				JWT: "f00b4r",
			},
		},
	}

	is := is.New(t)

	for _, tt := range tables {
		t.Run(tt.name, func(t *testing.T) {
			d := t.TempDir()

			config.Cfg.CacheDir = d

			t.Cleanup(func() { config.Cfg.CacheDir = "" })

			err := os.WriteFile(path.Join(d, "cache.json"), tt.prefill, 0o600)
			is.NoErr(err)

			err = cache.Save(tt.input)
			is.NoErr(err)

			l, err := cache.Load()
			is.NoErr(err)

			is.Equal(l.JWT, tt.expected.JWT)
			is.True(time.Time(l.Expires).Equal(time.Time(tt.expected.Expires)))
		})
	}
}
