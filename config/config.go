package config

import (
	"fmt"
	"os"

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

// LoadConfig loads the configuration using Viper and conditionally from config.properties
func LoadConfig() (*Config, error) {
	profile := os.Getenv("PROFILE")
	var config Config
	config.Profile = profile

	// 자동으로 환경 변수 로드
	viper.AutomaticEnv()

	if profile != "prod" {
		// 비프로덕션 환경에서는 config.properties 파일을 로드
		viper.SetConfigName("config")
		viper.SetConfigType("properties")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// 기본 서버 포트 설정
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}

	// 기본 SSLMode 설정
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	return &config, nil
}
