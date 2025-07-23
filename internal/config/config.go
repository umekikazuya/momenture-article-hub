package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// アプリケーションの全体設定を保持する。
type Config struct {
	AppEnv   string `mapstructure:"APP_ENV"`
	Database DatabaseConfig
}

// データベース接続設定を保持する。
type DatabaseConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_EXTERNAL_PORT"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	Name     string `mapstructure:"POSTGRES_DB"`
}

func LoadConfig(envFilePath string) (*Config, error) {
	viper.SetConfigFile(envFilePath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config file not found: %s. Relying on environment variables.\n", envFilePath)
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Config Structへのマッピング
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 環境変数からの上書き
	if appEnv := os.Getenv("APP_ENV"); appEnv != "" {
		config.AppEnv = appEnv
	}
	if dbHost := os.Getenv("POSTGRES_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("POSTGRES_EXTERNAL_PORT"); dbPort != "" {
		config.Database.Port = dbPort
	}
	if dbUser := os.Getenv("POSTGRES_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPassword := os.Getenv("POSTGRES_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if dbName := os.Getenv("POSTGRES_DB"); dbName != "" {
		config.Database.Name = dbName
	}

	if config.Database.Host == "" || config.Database.Port == "" || config.Database.User == "" || config.Database.Password == "" || config.Database.Name == "" {
		return nil, fmt.Errorf("database connection parameters are incomplete in config")
	}

	// データベースのPortを検証
	validPort, err := strconv.Atoi(config.Database.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid database port: %s", config.Database.Port)
	}
	if validPort <= 0 || validPort > 65535 {
		return nil, fmt.Errorf("database port must be between 1 and 65535: %d", validPort)
	}

	return &config, nil
}
