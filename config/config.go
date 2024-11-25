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
		Host     string `mapstructure:"HOST"`
		Name     string `mapstructure:"NAME"`
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
	fmt.Println("Profile:", profile) // 디버그 로그 추가
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
		// 비프로덕션 환경: Viper 사용
		// 자동으로 환경 변수 로드
		viper.AutomaticEnv()

		// 환경 변수의 '_'를 '.'으로 대체하여 Viper가 구조체 필드에 매핑할 수 있도록 함
		replacer := strings.NewReplacer("_", ".")
		viper.SetEnvKeyReplacer(replacer)

		// 비프로덕션 환경에서는 config.properties 파일을 로드
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

		fmt.Println("Host 2:", config.Database.Host) // 디버그 로그
	}

	fmt.Println("Host 3:", config.Database.Host) // 디버그 로그
	fmt.Printf("Loaded Config: %+v\n", config)   // 추가 디버그 로그

	// 기본 서버 포트 설정 (프로덕션과 비프로덕션 모두 적용)
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}

	// 기본 SSLMode 설정 (프로덕션과 비프로덕션 모두 적용)
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	// 필요한 환경 변수 검증 (프로덕션과 비프로덕션 모두 적용)
	if config.Database.Host == "" || config.Database.Name == "" || config.Database.User == "" || config.Database.Password == "" || config.Database.Port == "" {
		fmt.Printf("One or more required environment variables are missing: Host=%s, Name=%s, User=%s, Password=%s, Port=%s\n",
			config.Database.Host, config.Database.Name, config.Database.User, config.Database.Password, config.Database.Port)
		return nil, fmt.Errorf("one or more required environment variables are missing")
	}

	return &config, nil
}
