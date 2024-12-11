package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env              string     `yaml:"env" env-default:"local"`
	StoragePath      string     `yaml:"storage_path" env-required:"true"`
	GRPC             GRPCConfig `yaml:"grpc"`
	MigrationsPath   string
	TokenTTL         time.Duration          `yaml:"token_ttl" env-default:"1h"`
	Storage          DBConfig               `yaml:"storage"`
	StorageProcedure StorageProcedureConfig `yaml:"storage_procedure"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type DBConfig struct {
	Server   string `yaml:"server" env-required:"true"`
	Host     string `yaml:"host"`
	Path     string `yaml:"path"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database" env-required:"true"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
}

type StorageProcedureConfig struct {
	NameSP    string `yaml:"name_sp"`
	NameParam string `yaml:"name_param"`
	TvpType   string `yaml:"tvp_type"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
