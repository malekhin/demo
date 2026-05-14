package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Db         DbConfig         `mapstructure:"DB"`
	MonolithDb DbConfig         `mapstructure:"MONOLITH_DB"`
	Clickhouse ClickhouseConfig `mapstructure:"CLICKHOUSE"`
	Logger     LoggerConfig     `mapstructure:"LOGGER"`
	Redis      RedisConfig      `mapstructure:"REDIS"`
	AuthGate   AuthGateConfig   `mapstructure:"AUTHGATE"`
	Cron       CronConfig       `mapstructure:"CRON"`
	Metrics    MetricsConfig    `mapstructure:"METRICS"`
	App        AppConfig        `mapstructure:"APP"`
}

type DbConfig struct {
	Driver        string `mapstructure:"DRIVER" validate:"oneof=pgx mysql"`
	Host          string `mapstructure:"HOST"`
	Port          int    `mapstructure:"PORT"`
	User          string `mapstructure:"USER"`
	Password      string `mapstructure:"PASSWORD"`
	Name          string `mapstructure:"NAME"`
	QueryExecMode string `mapstructure:"QUERY_EXEC_MODE"`
	SslMode       string `mapstructure:"SSL_MODE"`
}

type ClickhouseConfig struct {
	Host         string        `mapstructure:"HOST"`
	Port         int           `mapstructure:"PORT"`
	User         string        `mapstructure:"USER"`
	Password     string        `mapstructure:"PASSWORD"`
	Name         string        `mapstructure:"NAME"`
	BatchWaiting time.Duration `mapstructure:"BATCH_WAITING"`
}

type LoggerConfig struct {
	Encoding string `mapstructure:"ENCODING"`
	Level    string `mapstructure:"LEVEL"`
}

type RedisConfig struct {
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Prefix   string `mapstructure:"PREFIX"`
}

type AuthGateConfig struct {
	Host  string        `mapstructure:"HOST"`
	Cache time.Duration `mapstructure:"CACHE"`
}

type CronConfig struct {
	ImportSubagent BaseLimitCronConfig `mapstructure:"IMPORT_SUBAGENT"`
	ActiveSubagent BaseCronConfig      `mapstructure:"ACTIVE_SUBAGENT"`
	Archive        BaseCronConfig      `mapstructure:"ARCHIVE"`
	Start          BaseCronConfig      `mapstructure:"START"`
	HistoryCalc    BaseLimitCronConfig `mapstructure:"HISTORY_CALC"`
}

type BaseLimitCronConfig struct {
	StartTime string `mapstructure:"START_TIME"`
	Limit     int    `mapstructure:"LIMIT"`
	Count     int    `mapstructure:"COUNT"`
	BatchSize int    `mapstructure:"BATCH_SIZE"`
}

type BaseCronConfig struct {
	StartTime string `mapstructure:"START_TIME" validate:"required"`
}

type MetricsConfig struct {
	Prefix string `mapstructure:"PREFIX"`
}

type AppConfig struct {
	SystemAgents []int         `mapstructure:"SYSTEM_AGENTS"`
	CalcCache    time.Duration `mapstructure:"CALC_CACHE"`
}

func New() (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()

	if cfgPath := os.Getenv("CONFIG_PATH"); cfgPath != "" {
		v.AddConfigPath(cfgPath)
	} else {
		v.AddConfigPath(".")
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigType("yaml")

	v.SetConfigName(fmt.Sprintf("config.%s", GetEnvironment()))

	if err := v.ReadInConfig(); err != nil {
		return &Config{}, fmt.Errorf("can't read config: %w", err)
	}

	cfg := Config{}
	err := v.Unmarshal(&cfg)
	if err != nil {
		return &Config{}, fmt.Errorf("can't unmarshal config: %w", err)
	}

	err = validator.New().Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}
