package config

import (
	"encoding"
	"fmt"
	"os"
)

var Cfg Config //nolint:gochecknoglobals

type Config struct {
	Email        string       `env:"EMAIL,required"`
	Password     string       `env:"PASSWORD,expand"`
	PasswordFile PasswordFile `env:"PASSWORD_FILE,expand"`
	CacheDir     string       `env:"CACHE_DIR,expand"     envDefault:"/var/cache/glucose_exporter"`
	Debug        bool         `env:"DEBUG"`
}

func (c *Config) GetPassword() (string, error) {
	if c.PasswordFile != "" && c.Password != "" {
		return "", ErrTooManyPasswords
	}

	if c.Password != "" {
		return c.Password, nil
	}

	if c.PasswordFile != "" {
		return string(c.PasswordFile), nil
	}

	return "", ErrMissingPassword
}

type PasswordFile string

var _ encoding.TextUnmarshaler = (*PasswordFile)(nil)

func (pf *PasswordFile) UnmarshalText(text []byte) error {
	b, err := os.ReadFile(string(text))
	if err != nil {
		return fmt.Errorf("reading password file: %w", err)
	}

	*pf = PasswordFile(string(b))

	return nil
}
