package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddr  string
	PostgresDSN string
}

func New(path string) (*Config, error) {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	flag.StringVar(&cfg.ServerAddr, "server_addr", os.Getenv("SERVER_ADDR"), "provide http server address")
	flag.StringVar(&cfg.PostgresDSN, "postgres_dsn", os.Getenv("POSTGRES_DSN"), "provide postgres database dsn")

	flag.Parse()

	strCfg, _ := json.MarshalIndent(&cfg, "", "  ")
	log.Println("Configuration:", string(strCfg))
	return &cfg, nil
}
