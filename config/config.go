package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Server          Server          `mapstructure:"server"`
		DatabasePostgre DatabasePostgre `mapstructure:"db_postgre"`
		DatabaseMySQL   DatabaseMySQL   `mapstructure:"db_mysql"`
		Redis           Redis           `mapstructure:"redis"`
		Logger          Logger          `mapstructure:"logger"`
		TelegramConfig  TelegramConfig  `mapstructure:"telegram"`
		JwtConfig       JwtConfig       `mapstructure:"jwt"`
	}

	Server struct {
		Port string `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	}

	DatabasePostgre struct {
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		User            string `mapstructure:"user"`
		Password        string `mapstructure:"password"`
		DBName          string `mapstructure:"dbname"`
		MaxIdleConns    int    `mapstructure:"maxIdleConns"`
		MaxOpenConns    int    `mapstructure:"maxOpenConns"`
		ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
	}

	DatabaseMySQL struct {
		Host            string `mapstructure:"host"`
		Port            int    `mapstructure:"port"`
		User            string `mapstructure:"user"`
		Password        string `mapstructure:"password"`
		DBName          string `mapstructure:"dbname"`
		MaxIdleConns    int    `mapstructure:"maxIdleConns"`
		MaxOpenConns    int    `mapstructure:"maxOpenConns"`
		ConnMaxLifetime int    `mapstructure:"connMaxLifetime"`
	}

	DatabaseMongo struct {
		URI             string `mapstructure:"uri"`
		Database        string `mapstructure:"database"`
		MaxPoolSize     int    `mapstructure:"maxPoolSize"`
		MinPoolSize     int    `mapstructure:"minPoolSize"`
		ConnMaxIdleTime int    `mapstructure:"connMaxIdleTime"`
	}

	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
		PoolSize int    `mapstructure:"poolSize"`
	}

	Logger struct {
		LogLevel   string            `mapstructure:"log_level"`
		FileLog    string            `mapstructure:"file_log"`
		MaxSize    int               `mapstructure:"max_size"`
		MaxBackups int               `mapstructure:"max_backups"`
		MaxAge     int               `mapstructure:"max_age"`
		Compress   bool              `mapstructure:"compress"`
		Channels   map[string]string `mapstructure:"channels"`
	}

	TelegramConfig struct {
		BotToken string `mapstructure:"bot_token"`
		ChatID   int64  `mapstructure:"chat_id"`
	}

	JwtConfig struct {
		Secret         string `mapstructure:"secret"`
		AccessTokenTtl string `mapstructure:"access_token_ttl"`
	}
)

var Cfg *Config

func Load(configDir string, configName string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName(configName)
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}
	return nil
}
