package config

import (
	"fmt"
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
	// 環境変数の自動読み込みを有効化
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// .envファイルの読み込み設定
	if envFilePath != "" {
		viper.SetConfigFile(envFilePath)
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Printf("Config file not found: %s. Relying on environment variables.\n", envFilePath)
			} else {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
		}
	}

	// 環境変数から設定を構築
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// ネストした構造体を個別にUnmarshal
	if err := viper.Unmarshal(&config.Database); err != nil {
		return nil, fmt.Errorf("failed to unmarshal database config: %w", err)
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
