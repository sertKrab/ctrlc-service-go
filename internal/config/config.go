package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv  string `mapstructure:"APP_ENV"`
	AppPort string `mapstructure:"APP_PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBSSLMode  string `mapstructure:"DB_SSL_MODE"`

	JWTSecret          string        `mapstructure:"JWT_SECRET"`
	AccessTokenExpiry  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRY"`
	RefreshTokenExpiry time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRY"`

	EnableRedis   bool   `mapstructure:"ENABLE_REDIS"`
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	CookieSecure   bool   `mapstructure:"COOKIE_SECURE"`
	CookieDomain   string `mapstructure:"COOKIE_DOMAIN"`
	CookieSameSite string `mapstructure:"COOKIE_SAME_SITE"`

	CORSOrigins string `mapstructure:"CORS_ORIGINS"`

	DBMaxOpenConns        int `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns        int `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBConnLifetimeMinutes int `mapstructure:"DB_CONN_LIFETIME_MINUTES"`

	DefaultLocale string `mapstructure:"DEFAULT_LOCALE"`

	LogLevel string `mapstructure:"LOG_LEVEL"`
}

func Load() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSL_MODE", "disable")
	viper.SetDefault("ACCESS_TOKEN_EXPIRY", "15m")
	viper.SetDefault("REFRESH_TOKEN_EXPIRY", "168h")
	viper.SetDefault("ENABLE_REDIS", false)
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("COOKIE_SECURE", false)
	viper.SetDefault("COOKIE_DOMAIN", "localhost")
	viper.SetDefault("COOKIE_SAME_SITE", "Lax")
	viper.SetDefault("CORS_ORIGINS", "http://localhost:3000")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_CONN_LIFETIME_MINUTES", 5)
	viper.SetDefault("DEFAULT_LOCALE", "th")
	viper.SetDefault("LOG_LEVEL", "info")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config unmarshal: %w", err)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET must be set")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME must be set")
	}
	return nil
}

func (c *Config) IsProduction() bool {
	return strings.ToLower(c.AppEnv) == "production"
}

func (c *Config) CookieSecureFlag() bool {
	if c.IsProduction() {
		return true
	}
	return c.CookieSecure
}
