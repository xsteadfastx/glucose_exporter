package config

var Cfg Config //nolint:gochecknoglobals

type Config struct {
	Email        string `env:"EMAIL,required"`
	Password     string `env:"PASSWORD,expand"`
	PasswordFile string `env:"PASSWORD_FILE,expand"`
	CacheDir     string `env:"CACHE_DIR,expand"     envDefault:"/var/cache/glucose_exporter"`
	Debug        bool   `env:"DEBUG"`
}
