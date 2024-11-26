package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the configuration values
type Config struct {
	Profile  string `mapstructure:"PROFILE"`
	Database struct {
		Host     string
		Name     string
		User     string `mapstructure:"USER"`
		Password string `mapstructure:"PASSWORD"`
		Port     string `mapstructure:"PORT"`
		SSLMode  string `mapstructure:"SSLMODE"`
	} `mapstructure:"DATABASE"`
	Server struct {
		Port string `mapstructure:"PORT"`
	} `mapstructure:"SERVER"`
}

func LoadConfig() (*Config, error) {
	profile := os.Getenv("PROFILE")
	fmt.Println("Profile:", profile)
	var config Config
	config.Profile = profile

	if profile == "prod" {
		// 프로덕션 환경: 환경 변수 직접 읽기
		config.Database.Host = os.Getenv("DATABASE_HOST")
		config.Database.Name = os.Getenv("DATABASE_NAME")
		config.Database.User = os.Getenv("DATABASE_USER")
		config.Database.Password = os.Getenv("DATABASE_PASSWORD")
		config.Database.Port = os.Getenv("DATABASE_PORT")
		config.Database.SSLMode = os.Getenv("DATABASE_SSLMODE")
		config.Server.Port = os.Getenv("SERVER_PORT") // 필요 시 추가

		fmt.Println("Host 1:", config.Database.Host) // 디버그 로그

	} else {
		viper.AutomaticEnv()

		replacer := strings.NewReplacer("_", ".")
		viper.SetEnvKeyReplacer(replacer)

		viper.SetConfigName("config")
		viper.SetConfigType("properties")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		err := viper.Unmarshal(&config)
		if err != nil {
			return nil, fmt.Errorf("unable to decode config into struct: %w", err)
		}
	}

	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}

	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	if config.Database.Host == "" || config.Database.Name == "" || config.Database.User == "" || config.Database.Password == "" || config.Database.Port == "" {
		fmt.Printf("One or more required environment variables are missing: Host=%s, Name=%s, User=%s, Password=%s, Port=%s\n",
			config.Database.Host, config.Database.Name, config.Database.User, config.Database.Password, config.Database.Port)
		return nil, fmt.Errorf("one or more required environment variables are missing")
	}

	return &config, nil
}
