package config

import (
	"balancer/internal/model"
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

type defaultLimits struct {
	Capacity   float64 `yaml:"capacity"`
	RatePerSec float64 `yaml:"rate_per_sec"`
}

type backendConfig struct {
	BackendURL string `yaml:"backend_url"`
	Config     struct {
		Health struct {
			Method string `yaml:"method"`
			URL    string `yaml:"url"`
		} `yaml:"health"`
	} `yaml:"config"`
}

type serverConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type dbConfig struct {
	DbHost     string `yaml:"db_host"`
	DbPort     string `yaml:"db_port"`
	DbName     string `yaml:"db_name"`
	DbUser     string `yaml:"db_user"`
	DbPassword string `yaml:"db_password"`
	DbSSL      string `yaml:"db_sslmode"`
}

type MainConfig struct {
	TickerRateSec uint64           `yaml:"ticker_rate_sec"`
	DefaultLimits defaultLimits    `yaml:"default_limits"`
	BackendConfig []*backendConfig `yaml:"backend_list"`
	ServerConfig  serverConfig     `yaml:"server"`
	DbConfig      dbConfig         `yaml:"db"`
}

func InitMainConfig(configPath string) (*MainConfig, error) {
	var cfg MainConfig
	if _, err := os.Stat(configPath); err != nil {
		return nil, model.ErrParseConfig
	}
	rowConfig, err := os.ReadFile(configPath)
	if err != nil {
		return nil, model.ErrParseConfig
	}
	err = yaml.Unmarshal(rowConfig, &cfg)
	if err != nil {
		return nil, model.ErrParseConfig
	}
	return &cfg, nil
}

func (cfg *MainConfig) LoadTickerRateSec() uint64 {
	return cfg.TickerRateSec
}

func (cfg *MainConfig) LoadDefaultLimits() model.DefaultClientLimits {
	return model.DefaultClientLimits{
		Capacity:   cfg.DefaultLimits.Capacity,
		RatePerSec: cfg.DefaultLimits.RatePerSec,
	}
}

func (cfg *MainConfig) LoadBackendConfig() []model.BackendServerSettings {
	settings := make([]model.BackendServerSettings, 0, len(cfg.BackendConfig))
	for _, b := range cfg.BackendConfig {
		settings = append(settings, model.BackendServerSettings{
			BckndUrl: b.BackendURL,
			Method:   b.Config.Health.Method,
			HelthUrl: b.Config.Health.URL,
		})
	}
	return settings
}

func (cfg *MainConfig) LoadServerAddress() string {
	return net.JoinHostPort(cfg.ServerConfig.Host, cfg.ServerConfig.Port)
}

func (cfg *MainConfig) LoadDbConfig() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.DbConfig.DbHost,
		cfg.DbConfig.DbPort,
		cfg.DbConfig.DbName,
		cfg.DbConfig.DbUser,
		cfg.DbConfig.DbPassword,
		cfg.DbConfig.DbSSL,
	)
}
