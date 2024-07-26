package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Postgres   Postgres   `yaml:"postgres"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Redis      Redis      `yaml:"redis"`
	Session    Session    `yaml:"session"`
	Cookie     Cookie     `yaml:"cookie"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DbName   string `yaml:"db_name"`
	Password string `yaml:"password"`
}

type HTTPServer struct {
	JWTSecretKey string        `yaml:"jwt_secret_key"`
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type Redis struct {
	Address  string `yaml:"address"`
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
}

type Session struct {
	Prefix string `yaml:"prefix"`
	Name   string `yaml:"name"`
	Expire int    `yaml:"expire"`
}

type Cookie struct {
	Name     string `yaml:"name"`
	MaxAge   int    `yaml:"max_age"`
	Secure   bool   `yaml:"secure"`
	HTTPOnly bool   `yaml:"http_only"`
}

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	cfg := &Config{}
	cleanenv.ReadConfig(configPath, cfg)

	return cfg
}
