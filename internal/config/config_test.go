package config_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/matryer/is"
	"go.xsfx.dev/glucose_exporter/internal/config"
)

func TestPasswordFile(t *testing.T) {
	is := is.New(t)

	dir := t.TempDir()

	pfPath := path.Join(dir, "foo.txt")

	err := os.WriteFile(pfPath, []byte("f00b4r"), 0o600)
	is.NoErr(err)

	var pf config.PasswordFile

	err = pf.UnmarshalText([]byte(pfPath))
	is.NoErr(err)

	is.Equal(string(pf), "f00b4r")
}

func Test(t *testing.T) {
	tables := []struct {
		name     string
		cfg      config.Config
		expected string
		err      error
	}{
		{
			name:     "00",
			cfg:      config.Config{},
			expected: "",
			err:      config.ErrMissingPassword,
		},
		{
			name: "01",
			cfg: config.Config{
				Password:     "foo",
				PasswordFile: "bar",
			},
			expected: "",
			err:      config.ErrTooManyPasswords,
		},
		{
			name: "02",
			cfg: config.Config{
				Password: "foo",
			},
			expected: "foo",
			err:      nil,
		},
		{
			name: "03",
			cfg: config.Config{
				PasswordFile: "foo",
			},
			expected: "foo",
			err:      nil,
		},
	}

	is := is.New(t)

	for _, tt := range tables {
		t.Run(tt.name, func(_ *testing.T) {
			pass, err := tt.cfg.GetPassword()
			if tt.err == nil {
				is.NoErr(err)
				is.Equal(pass, tt.expected)
			} else {
				is.True(errors.Is(err, tt.err))
			}
		})
	}
}
