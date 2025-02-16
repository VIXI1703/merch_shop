package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	DB   DB   `mapstructure:"database"`
	HTTP HTTP `mapstructure:"http"`
	JWT  JWT  `mapstructure:"jwt"`
}

type JWT struct {
	SigningKey string        `mapstructure:"signing_key"`
	Duration   time.Duration `mapstructure:"duration"`
}

type DB struct {
	User           string        `mapstructure:"user"`
	Password       string        `mapstructure:"password"`
	Name           string        `mapstructure:"name"`
	Host           string        `mapstructure:"host"`
	Port           string        `mapstructure:"port"`
	DBMaxOpenConns int           `mapstructure:"max_open_conns"`
	DBMaxIdleConns int           `mapstructure:"max_idle_conns"`
	DBConnMaxLife  time.Duration `mapstructure:"conn_max_life"`
}

type HTTP struct {
	Port string `mapstructure:"port"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")

	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_life", time.Hour)
	viper.SetDefault("jwt.duration", time.Hour*24)

	viper.AutomaticEnv()
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.conn_max_life", "DB_CONN_MAX_LIFE")
	viper.BindEnv("database.max_idle_conns", "DB_MAX_IDLE_CONNS")
	viper.BindEnv("database.max_open_conns", "DB_MAX_OPEN_CONNS")
	viper.BindEnv("jwt.signing_key", "JWT_SIGNING_KEY")
	viper.BindEnv("jwt.duration", "JWT_DURATION")
	viper.BindEnv("http.port", "HTTP_PORT")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config, func(decoder *mapstructure.DecoderConfig) { decoder.ErrorUnset = true }); err != nil {
		return Config{}, fmt.Errorf("unable to decode into struct, %w", err)
	}

	return config, nil
}
