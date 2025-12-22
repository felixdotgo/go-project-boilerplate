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
	Secret          string `mapstructure:"jwt_secret"`
	Expire          string `mapstructure:"jwt_expire"`
	Issuer          string `mapstructure:"jwt_issuer"`
	AccessTokenTTL  string `mapstructure:"jwt_access_token_ttl"`  // Default: 1h
	RefreshTokenTTL string `mapstructure:"jwt_refresh_token_ttl"` // Default: 720h (30 days)
	EncryptionKey   string `mapstructure:"oauth_encryption_key"`  // For encrypting OAuth provider tokens
}

// OAuth contains all configurations that related to OAuth
type OAuth struct {
	RedirectURL string              `mapstructure:"oauth_redirect_url"`
	Google      GoogleOAuthConfig   `mapstructure:",squash"`
	Facebook    FacebookOAuthConfig `mapstructure:",squash"`
	GitHub      GitHubOAuthConfig   `mapstructure:",squash"`
}

type OAuther interface {
	IsEnabled() bool
	GetScopes() []string
}

// GoogleOAuthConfig contains configs for Google OAuth provider
type GoogleOAuthConfig struct {
	ClientID     string   `mapstructure:"oauth_google_client_id"`
	ClientSecret string   `mapstructure:"oauth_google_client_secret"`
	RedirectURL  string   `mapstructure:"oauth_google_redirect_url"`
	Scopes       string `mapstructure:"oauth_google_scopes"`
	Enabled      string   `mapstructure:"oauth_google_enabled"`
}

func (o *GoogleOAuthConfig) IsEnabled() bool {
	return strings.ToLower(o.Enabled) == "true"
}

func (o *GoogleOAuthConfig) GetScopes() []string {
	if o.Scopes == "" {
		return []string{}
	}
	return strings.Split(o.Scopes, " ")
}

// FacebookOAuthConfig contains configs for Facebook OAuth provider
type FacebookOAuthConfig struct {
	ClientID     string   `mapstructure:"oauth_facebook_client_id"`
	ClientSecret string   `mapstructure:"oauth_facebook_client_secret"`
	RedirectURL  string   `mapstructure:"oauth_facebook_redirect_url"`
	Scopes       string `mapstructure:"oauth_facebook_scopes"`
	Enabled      string   `mapstructure:"oauth_facebook_enabled"`
}

func (o *FacebookOAuthConfig) IsEnabled() bool {
	return strings.ToLower(o.Enabled) == "true"
}

func (o *FacebookOAuthConfig) GetScopes() []string {
	if o.Scopes == "" {
		return []string{}
	}
	return strings.Split(o.Scopes, " ")
}

// GitHubOAuthConfig contains configs for GitHub OAuth provider
type GitHubOAuthConfig struct {
	ClientID     string   `mapstructure:"oauth_github_client_id"`
	ClientSecret string   `mapstructure:"oauth_github_client_secret"`
	RedirectURL  string   `mapstructure:"oauth_github_redirect_url"`
	Scopes       string `mapstructure:"oauth_github_scopes"`
	Enabled      string   `mapstructure:"oauth_github_enabled"`
}

func (o *GitHubOAuthConfig) IsEnabled() bool {
	return strings.ToLower(o.Enabled) == "true"
}

func (o *GitHubOAuthConfig) GetScopes() []string {
	if o.Scopes == "" {
		return []string{}
	}
	return strings.Split(o.Scopes, " ")
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
	v.AddConfigPath("./cmd/svc-auth")
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
