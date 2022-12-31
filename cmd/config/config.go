package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port            string
	PostgresDSN     string
	AccessTokenKey  string
	RefreshTokenKey string
	AccessTokenExp  int
	RefreshTokenExp int
}

func New() *Config {
	var cfg Config

	portEnv := os.Getenv("PORT")
	postgreDSNEnv := os.Getenv("POSTGRES_DSN")
	accessTokenKeyEnv := os.Getenv("ACCESS_TOKEN_KEY")
	refreshTokenKeyEnv := os.Getenv("REFRESH_TOKEN_KEY")
	accessTokenExpEnv, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXP"))
	refreshTokenExpEnv, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXP"))

	flag.StringVar(&cfg.Port, "port", ":"+portEnv, "provide http server port address")
	flag.StringVar(&cfg.PostgresDSN, "postgres_dsn", postgreDSNEnv, "provide postgres database dsn")
	flag.StringVar(&cfg.AccessTokenKey, "access_token_key", accessTokenKeyEnv, "provide access token secret key for jwt")
	flag.StringVar(&cfg.RefreshTokenKey, "refresh_token_key", refreshTokenKeyEnv, "provide refresh token secret key for jwt")
	flag.IntVar(&cfg.AccessTokenExp, "access_token_exp", accessTokenExpEnv, "provide access token expiration time in seconds")
	flag.IntVar(&cfg.RefreshTokenExp, "refresh_token_exp", refreshTokenExpEnv, "provide refresh token expiration time in seconds")

	flag.Parse()

	strCfg, _ := json.MarshalIndent(&cfg, "", "  ")
	log.Println("Configuration:", string(strCfg))
	return &cfg
}
