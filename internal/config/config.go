// Package config
/*
Переопределение переменных осуществляется file -> cmd -> env,
то есть переменные окружения ENV_* переопределят всё остальное
*/
package config

import (
	"github.com/cepmap/otus-system-monitoring/internal/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
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
		Cpu         bool  `mapstructure:"cpu" env:"STATS_CPU"`
		DiskInfo    bool  `mapstructure:"disk_info" env:"STATS_DISK_INFO"`
		DiskLoad    bool  `mapstructure:"disk_load" env:"STATS_DISK_LOAD"`
		NetStat     bool  `mapstructure:"net_stat" env:"STATS_NET_STAT"`
		TopTalkers  bool  `mapstructure:"top_talkers" env:"STATS_TOP_TALKERS"`
	} `mapstructure:"stats"`
}

var DaemonConfig *Config

func init() {

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
	}

	if err := viper.BindPFlag("log.level", pflag.Lookup("loglevel")); err != nil {
		logger.Error(err.Error())
	}
	if err := viper.BindPFlag("server.host", pflag.Lookup("host")); err != nil {
		logger.Error(err.Error())
	}
	if err := viper.BindPFlag("server.port", pflag.Lookup("port")); err != nil {
		logger.Error(err.Error())
	}

	if err := viper.Unmarshal(DaemonConfig); err != nil {
		logger.Error(err.Error())
	}
}

func initSettings() Config {
	return Config{
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
			Cpu         bool  `mapstructure:"cpu" env:"STATS_CPU"`
			DiskInfo    bool  `mapstructure:"disk_info" env:"STATS_DISK_INFO"`
			DiskLoad    bool  `mapstructure:"disk_load" env:"STATS_DISK_LOAD"`
			NetStat     bool  `mapstructure:"net_stat" env:"STATS_NET_STAT"`
			TopTalkers  bool  `mapstructure:"top_talkers" env:"STATS_TOP_TALKERS"`
		}{LoadAverage: true, Cpu: false, DiskInfo: false, NetStat: false, TopTalkers: false},
	}
}
