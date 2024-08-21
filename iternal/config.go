package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-default:"./storage/storage.db"`
	HttpServer  HTTPServer
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8000"`
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        int           `yaml:"port" env-default:"8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"25s"`
	User        string        `yaml:"user"  env-default:"admin"`
	Password    string        `yaml:"password"  env-default:"admin" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad() *Config {
	cfgPath := "./config/config.yaml"

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH is not exist: %s", cfgPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("config error: %s", err)
	}

	return &cfg

}
