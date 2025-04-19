package configs

import (
	"errors"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CacheType string

// Config stores all configuration of the application.
type Config struct {
	Token    TokenConfig
	Database DatabaseConfig
	Cache    CacheConfig
	HTTP     HTTPConfig
	Log      Log
}

// TokenConfig struct for token configuration
type TokenConfig struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	TokenSymmetricKey    string
}

// DatabaseConfig struct for database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// CacheConfig struct for cache configuration
type CacheConfig struct {
	Type     CacheType
	Host     string
	Port     int
	Username string
	Password string
}

// HTTPConfig struct for HTTP server configuration
type HTTPConfig struct {
	Address string
}

// HTTPConfig struct for HTTP server configuration
type Log struct {
	Level string
}

func LoadConfig() (config Config, err error) {
	// Define the flag
	configFileFlag := pflag.StringP("config", "f", "", "Path to the configuration file")
	// Parse the flags
	pflag.Parse()

	// Check if the flag was provided
	if *configFileFlag == "" {
		err = errors.New("flag empty")
		return
	}

	viper.SetConfigFile(*configFileFlag)
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
