package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// LOADS CONFIG
func Load(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := CreateDefaultConfig(configPath); err != nil {
			log.Printf("CONFIG: Failed to create default config file: %v", err)
			return nil, errors.New("failed to create default config file")
		}
		log.Printf("CONFIG: Created default config file at %s", configPath)
		log.Printf("CONFIG: Please review and update the configuration before running again")
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("CONFIG: Failed to read config file: %v", err)
		return nil, errors.New("failed to read config file")
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		log.Printf("CONFIG: Failed to parse config file: %v", err)
		return nil, errors.New("failed to parse config file")
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// CREATES DEFAULT CONFIG IF MISSING
func CreateDefaultConfig(configPath string) error {
	cfg := DefaultConfig()

	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshalling: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing: %w", err)
	}

	return nil
}

// VALIDATES CONFIG
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return errors.New("server port is required")
	}
	if c.Database.Path == "" {
		return errors.New("database path is required")
	}
	if c.Security.JWTSecret == "" {
		return errors.New("JWT secret is required")
	}
	if c.Security.JWTExpiryHours <= 0 {
		return errors.New("JWT expiry hours must be positive")
	}
	return nil
}

// RETURNS DB STRING FROM CONFIG
func (d *DatabaseConfig) GetConnectionString() (string, error) {
	switch strings.ToLower(d.Driver) {
	case "sqlite3":
		if d.Path == "" {
			return "", fmt.Errorf("database path is required for sqlite3")
		}
		return d.Path, nil
	case "mysql":
		if d.Host == "" || d.Port == "" || d.User == "" || d.DBName == "" {
			return "", fmt.Errorf("host, port, user, and dbname are required for mysql")
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", d.User, d.Password, d.Host, d.Port, d.DBName), nil
	case "postgres":
		if d.Host == "" || d.Port == "" || d.User == "" || d.DBName == "" {
			return "", fmt.Errorf("host, port, user, and dbname are required for postgres")
		}
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode), nil
	default:
		return "", fmt.Errorf("unsupported database driver: %s (supported: sqlite3, mysql, postgres)", d.Driver)
	}
}
