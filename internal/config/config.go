/*
Переопределение переменных осуществляется file -> cmd -> env,
то есть переменные окружения ENV_* переопределят всё остальное
*/
package config

import (
	"fmt"
	"strings"

	"github.com/cepmap/otus-system-monitoring/internal/logger"
	"github.com/cepmap/otus-system-monitoring/internal/tools"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Log struct {
		Level string `mapstructure:"level" env:"LOG_LEVEL"`
	} `mapstructure:"log"`
	Server struct {
		Host string `mapstructure:"host" env:"SERVER_HOST"`
		Port string `mapstructure:"port" env:"SERVER_PORT"`
	} `mapstructure:"server"`
	Stats struct {
		Limit       int64 `mapstructure:"limit" env:"STATS_LIMIT"`
		LoadAverage bool  `mapstructure:"load_average" env:"STATS_LOAD_AVERAGE"`
		Cpu         bool  `mapstructure:"CPU" env:"STATS_CPU"`
		DiskInfo    bool  `mapstructure:"disk_info" env:"STATS_DISK_INFO"`
		DiskLoad    bool  `mapstructure:"disk_load" env:"STATS_DISK_LOAD"`
	} `mapstructure:"stats"`
}

var DaemonConfig *Config

func InitConfig() error {
	initSettings := initSettings()
	DaemonConfig = &initSettings

	configFilePath := pflag.String("config", "./_configs/config.yaml", "Config file")
	pflag.String("loglevel", "info", "Log level")
	pflag.String("host", "0.0.0.0", "Server host")
	pflag.String("port", "8080", "Server port")
	pflag.Parse()

	viper.SetConfigFile(*configFilePath)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.BindPFlag("log.level", pflag.Lookup("loglevel")); err != nil {
		return fmt.Errorf("failed to bind log level flag: %w", err)
	}
	if err := viper.BindPFlag("server.host", pflag.Lookup("host")); err != nil {
		return fmt.Errorf("failed to bind host flag: %w", err)
	}
	if err := viper.BindPFlag("server.port", pflag.Lookup("port")); err != nil {
		return fmt.Errorf("failed to bind port flag: %w", err)
	}

	if err := viper.Unmarshal(DaemonConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	checkCommands(DaemonConfig)
	return nil
}

func initSettings() Config {
	config := Config{
		Log: struct {
			Level string `mapstructure:"level" env:"LOG_LEVEL"`
		}{Level: "DEBUG"},
		Server: struct {
			Host string `mapstructure:"host" env:"SERVER_HOST"`
			Port string `mapstructure:"port" env:"SERVER_PORT"`
		}{Host: "0.0.0.0", Port: "8080"},
		Stats: struct {
			Limit       int64 `mapstructure:"limit" env:"STATS_LIMIT"`
			LoadAverage bool  `mapstructure:"load_average" env:"STATS_LOAD_AVERAGE"`
			Cpu         bool  `mapstructure:"CPU" env:"STATS_CPU"`
			DiskInfo    bool  `mapstructure:"disk_info" env:"STATS_DISK_INFO"`
			DiskLoad    bool  `mapstructure:"disk_load" env:"STATS_DISK_LOAD"`
		}{LoadAverage: true, Cpu: false, DiskInfo: false, DiskLoad: false},
	}
	return config
}

func checkCommands(config *Config) {
	if config.Stats.Cpu {
		if err := tools.CheckCommand("iostat"); err != nil {
			logger.Error("command iostat not found, disabling cpu stats collection")
			config.Stats.Cpu = false
		}
	}
	if config.Stats.DiskLoad {
		if err := tools.CheckCommand("iostat"); err != nil {
			logger.Error("command iostat not found, disabling disk load stats collection")
			config.Stats.DiskLoad = false
		}
	}
}
