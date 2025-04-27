package config

import (
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

type serverConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type mainConfig struct {
	Server      serverConfig `yaml:"server"`
	BackendList []string     `yaml:"backend_list"`
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

func (cfg *mainConfig) GetBackendList() []string {
	return cfg.BackendList
}

func (cfg *mainConfig) GetServerAddress() string {
	return net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)
}
