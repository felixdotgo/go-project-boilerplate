package config

import (
	"fmt"
	"strings"

	"github.com/0x46656C6978/go-project-boilerplate/pkg/conv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const ENV_PRODUCTION = "production"

// DB contains configurations to make a connection to RDBMS
type DB struct {
	Host     string `mapstructure:"db_host"`
	User     string `mapstructure:"db_user"`
	Password string `mapstructure:"db_password"`
	DBName   string `mapstructure:"db_dbname"`
	Port     string `mapstructure:"db_port"`
}

// JWT is configurations for JSON Web Token
type JWT struct {
	Secret string `mapstructure:"jwt_secret"`
	Expire string `mapstructure:"jwt_expire"`
	Issuer string `mapstructure:"jwt_issuer"`
}

// OAuth contains all configurations that related to OAuth
type OAuth struct {
	Google OAuthGoogle `mapstructure:",squash"`
}

// OAuthGoogle contains configs that will be used to make an auth request to Google
type OAuthGoogle struct {
	ClientID     string `mapstructure:"oauth_google_client_id"`
	ClientSecret string `mapstructure:"oauth_google_client_secret"`
	RedirectURL  string `mapstructure:"oauth_redirect_url"`
}

// Config is a struct that contains all other configurations that will be defined later
type Config struct {
	EnvMode string `mapstructure:"env_mode"`
	Port    string `mapstructure:"port"`
	DB      DB     `mapstructure:",squash"`
	JWT     JWT    `mapstructure:",squash"`
	OAuth   OAuth  `mapstructure:",squash"`
}

// New returns new Config
func New() (*Config, error) {
	// solution: https://stackoverflow.com/a/63541140/9839165
	v := viper.New()
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)
	// Load config from .env file
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	v.AddConfigPath(".")
	v.AddConfigPath("./cmd/api")
	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s", err)
	}

	var result map[string]interface{}
	var cfg *Config

	if err := v.Unmarshal(&result); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %s", err)
	}

	if err := mapstructure.Decode(result, &cfg); err != nil {
		return nil, fmt.Errorf("error decoding config: %s", err)
	}

	return cfg, nil
}

// IsProduction check whether env mode is production or not
func (c *Config) IsProduction() bool {
	if c == nil {
		return false
	}
	return c.EnvMode == ENV_PRODUCTION
}

// GetPort return string value of Port as int
func (c *Config) GetPort() int {
	if c == nil {
		return 0
	}
	return conv.ToInt(c.Port)
}

// GetEnvMode return string value of EnvMode
func (c *Config) GetEnvMode() string {
	if c == nil {
		return ""
	}
	return c.EnvMode
}
