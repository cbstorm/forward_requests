package configs

import (
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var cfg *Config
var cfg_sync sync.Once

func GetConfig() *Config {
	if cfg == nil {
		cfg_sync.Do(func() {
			cfg = NewConfig()
		})
	}
	return cfg
}

func NewConfig() *Config {
	return &Config{}
}

type Config struct {
	ENV      string
	APP_PORT string
	SERVERS  []string
}

func (cfg *Config) Load() error {
	if os.Getenv("ENV") == "development" {
		env_path := ".env"
		err := godotenv.Load(env_path)
		if err != nil {
			return err
		}
	}
	cfg.ENV = os.Getenv("ENV")
	cfg.APP_PORT = os.Getenv("APP_PORT")
	server_strs := os.Getenv("SERVERS")
	cfg.SERVERS = strings.Split(server_strs, "|")
	return nil
}
