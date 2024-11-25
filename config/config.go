package config

import (
	"fmt"

	"github.com/go-ini/ini"
	"github.com/spf13/viper"
)

// Config holds the configuration values
type Config struct {
	Profile  string
	Database struct {
		Host     string
		Name     string
		User     string
		Password string
		Port     string
		SSLMode  string
	}
	Server struct {
		Port string
	}
}

// LoadConfig loads the configuration using Viper
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("properties")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Set default server port if not set
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}

	// Set default SSLMode if not set
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	return &config, nil
}

// LoadConfigFromIni loads configuration from config.properties for non-production
func LoadConfigFromIni(cfg *Config) error {
	iniFile, err := ini.Load("config.properties")
	if err != nil {
		return fmt.Errorf("failed to read config.properties: %w", err)
	}

	dbSection := iniFile.Section("database")
	cfg.Database.Host = dbSection.Key("DATABASE_HOST").String()
	cfg.Database.Name = dbSection.Key("DATABASE_NAME").String()
	cfg.Database.User = dbSection.Key("DATABASE_USER").String()
	cfg.Database.Password = dbSection.Key("DATABASE_PASSWORD").String()
	cfg.Database.Port = dbSection.Key("DATABASE_PORT").String()

	// Additional validation
	if cfg.Database.Host == "" || cfg.Database.Name == "" || cfg.Database.User == "" || cfg.Database.Password == "" || cfg.Database.Port == "" {
		return fmt.Errorf("one or more required configuration values are missing in config.properties")
	}

	return nil
}
