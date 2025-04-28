package config

import (
	"balancer/internal/model"
	"fmt"
	"log"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

type defoultLimits struct {
	Capacity   float64 `yaml:"capasity"`
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

type mainConfig struct {
	TickerRate    uint64          `yaml:"ticker_rate_inseconds"`
	DefoultLimits defoultLimits   `yaml:"defoult_limits"`
	BackendList   []backendConfig `yaml:"backend_list"`
	Server        serverConfig    `yaml:"server"`
	DbConfig      dbConfig        `yaml:"db"`
}

func InitMainConfig(configPath string) (mainConfig, error) {
	var cfg mainConfig
	if _, err := os.Stat(configPath); err != nil {
		return mainConfig{}, err
	}
	rowConfig, err := os.ReadFile(configPath)
	if err != nil {
		return mainConfig{}, err
	}
	err = yaml.Unmarshal(rowConfig, &cfg)
	if err != nil {
		return mainConfig{}, err
	}
	return cfg, nil
}

func (cfg *mainConfig) GetTickerRate() uint64 {
	return cfg.TickerRate
}

func (cfg *mainConfig) GetDefoulLimits() *model.DefoultUserLimits {
	return &model.DefoultUserLimits{
		Capacity:   cfg.DefoultLimits.Capacity,
		RatePerSec: cfg.DefoultLimits.RatePerSec,
	}
}

func (cfg *mainConfig) InitBackendList() []*model.BackendServer {
	backendList := make([]*model.BackendServer, len(cfg.BackendList))
	for i, b := range cfg.BackendList {
		if b.Config.Health.URL == "" {
			b.Config.Health.URL = "/health"
		}
		if b.Config.Health.Method == "" {
			b.Config.Health.Method = "GET"
		}
		backendList[i] = &model.BackendServer{
			BckndUrl: b.BackendURL,
			Method:   b.Config.Health.Method,
			HelthUrl: b.BackendURL + b.Config.Health.URL,
		}
		log.Println("Loaded backend:", b.BackendURL, "Method:", b.Config.Health.Method)
	}
	return backendList
}

func (cfg *mainConfig) GetServerAddress() string {
	return net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
}

func (cfg *mainConfig) DbConfigLoad() string {

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
