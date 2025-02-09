package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string
}

// env-default:"production"
type Config struct {
	Env string `yaml:"env" env:"ENV" env-required:"true" `
	StoragePath string
	HTTPServer

}

func MustLoad() *Config{
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the config file")

		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("config path is not set ")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist : %s", configPath)
	}
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatal("can not read config file : %s",err.Error() )
	}

	return &cfg;

}