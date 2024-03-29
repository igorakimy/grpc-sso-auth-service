package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const (
	configPathEnv = "CONFIG_PATH"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	TokenTTL time.Duration `yaml:"token_ttl"`
	Grpc     GRPCConfig    `yaml:"grpc"`
	Db       DbConfig      `yaml:"db"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DSN      string `yaml:"dsn"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path not specified")
	}

	return MustLoadFromPath(path)
}

func MustLoadFromPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv(configPathEnv)
	}

	return path
}
