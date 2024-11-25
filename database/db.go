package database

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/go-ini/ini"
	_ "github.com/lib/pq"
	"zim-iot-data-api/config"
)

// DB wraps the sql.DB to allow for future extensions
type DB struct {
	*sql.DB
	Logger *logrus.Logger
}

// NewDB initializes the database connection
func NewDB(cfg *config.Config, logger *logrus.Logger) (*DB, error) {
	var connStr string
	if cfg.Profile == "prod" {
		if cfg.Database.Host == "" || cfg.Database.Name == "" || cfg.Database.User == "" || cfg.Database.Password == "" || cfg.Database.Port == "" {
			return nil, fmt.Errorf("one or more required environment variables are missing")
		}
		connStr = fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=%s",
			cfg.Database.Host, cfg.Database.Name, cfg.Database.User, cfg.Database.Password, cfg.Database.Port, cfg.Database.SSLMode)
	} else {
		// Load from config.properties
		iniFile, err := ini.Load("config.properties")
		if err != nil {
			return nil, fmt.Errorf("failed to read config.properties: %w", err)
		}
		dbSection := iniFile.Section("database")
		dbHost := dbSection.Key("DATABASE_HOST").String()
		dbName := dbSection.Key("DATABASE_NAME").String()
		dbUser := dbSection.Key("DATABASE_USER").String()
		dbPassword := dbSection.Key("DATABASE_PASSWORD").String()
		dbPort := dbSection.Key("DATABASE_PORT").String()
		if dbHost == "" || dbName == "" || dbUser == "" || dbPassword == "" || dbPort == "" {
			return nil, fmt.Errorf("one or more required configuration values are missing in config.properties")
		}
		connStr = fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=%s",
			dbHost, dbName, dbUser, dbPassword, dbPort, cfg.Database.SSLMode)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping the database to verify connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("Connected to database successfully.")

	return &DB{
		DB:     db,
		Logger: logger,
	}, nil
}
