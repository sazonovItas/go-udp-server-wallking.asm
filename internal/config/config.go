package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env       string `yaml:"env"        env-default:"local"`
	UDPServer `       yaml:"udp_server"`
}

type UDPServer struct {
	Address     string        `yaml:"address"      env-default:"localhost"`
	Port        string        `yaml:"port"         env-default:"8829"`
	Timeout     time.Duration `yaml:"timeout"      env-default:"1s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func Load() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %s", err)
	}

	return &cfg, nil
}
