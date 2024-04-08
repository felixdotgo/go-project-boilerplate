package config

import (
	"os"
	"time"

	"github.com/0x46656C6978/go-project-boilerplate/pkg/conv"
)

const ENV_PRODUCTION = "production"

type Config struct {
	EnvMode string
	Port    int
	DB      struct {
		Host     string
		User     string
		Password string
		DBName   string
		Port     int
	}
	JWT struct {
		Secret string
		Expire time.Duration
		Issuer string
	}
	OAuth struct {
		Google struct {
			ClientID     string
			ClientSecret string
			RedirectURL	string
		}
	}
}

func New() *Config {
	cfg := &Config{}
	cfg.EnvMode = os.Getenv("ENV_MODE")
	cfg.Port = conv.ToInt(os.Getenv("PORT"))

	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.DBName = os.Getenv("DB_DBNAME")
	cfg.DB.Port = conv.ToInt(os.Getenv("DB_PORT"))

	cfg.JWT.Expire = time.Duration(conv.ToInt(os.Getenv("JWT_EXPIRE"))) * time.Second
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	cfg.JWT.Issuer = os.Getenv("JWT_ISSUER")

	cfg.OAuth.Google.ClientID = os.Getenv("OAUTH_GOOGLE_CLIENT_ID")
	cfg.OAuth.Google.ClientSecret = os.Getenv("OAUTH_GOOGLE_CLIENT_SECRET")
	cfg.OAuth.Google.RedirectURL = os.Getenv("OAUTH_GOOGLE_REDIRECT_URL")

	return cfg
}

func (c *Config) IsProduction() bool {
	return c.EnvMode == ENV_PRODUCTION
}
