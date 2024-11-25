package config

import (
	"fmt"
	"os"
	"strings"

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
	fmt.Println("Profile:", profile)
	var config Config
	config.Profile = profile

	// 자동으로 환경 변수 로드
	viper.AutomaticEnv()

	// 환경 변수의 '_'를 '.'으로 대체하여 Viper가 구조체 필드에 매핑할 수 있도록 함
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	fmt.Println("Host 1", config.Database.Host)

	if profile != "prod" {
		// 비프로덕션 환경에서는 config.properties 파일을 로드
		viper.SetConfigName("config")
		viper.SetConfigType("properties")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	fmt.Println("Host 2", config.Database.Host)

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	fmt.Println("Host 3", config.Database.Host)

	// 기본 서버 포트 설정
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}

	// 기본 SSLMode 설정
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	// 필요한 환경 변수 검증
	if config.Database.Host == "" || config.Database.Name == "" || config.Database.User == "" || config.Database.Password == "" || config.Database.Port == "" {
		fmt.Printf("One or more required environment variables are missing: Host=%s, Name=%s, User=%s, Password=%s, Port=%s\n",
			config.Database.Host, config.Database.Name, config.Database.User, config.Database.Password, config.Database.Port)
		return nil, fmt.Errorf("one or more required environment variables are missing")
	}

	return &config, nil
}
