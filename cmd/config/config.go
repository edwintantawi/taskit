package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"

	"github.com/edwintantawi/taskit/pkg/postgres"
)

type Config struct {
	Port                   string
	AllowedOrigin          string
	AccessTokenKey         string
	RefreshTokenKey        string
	AccessTokenExpiration  int
	RefreshTokenExpiration int
	AutoMigrate            bool
	Postgres               postgres.Config
}

func New() *Config {
	var config Config

	portEnv := os.Getenv("PORT")
	allowedOriginEnv := os.Getenv("ALLOWED_ORIGIN")
	accessTokenKeyEnv := os.Getenv("ACCESS_TOKEN_KEY")
	refreshTokenKeyEnv := os.Getenv("REFRESH_TOKEN_KEY")
	accessTokenExpirationEnv, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRATION"))
	refreshTokenExpirationEnv, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION"))
	autoMigrateEnv, _ := strconv.ParseBool(os.Getenv("AUTO_MIGRATE"))

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDB := os.Getenv("POSTGRES_DB")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresSSLModeEnv := os.Getenv("POSTGRES_SSLMODE")

	flag.StringVar(&config.Port, "port", portEnv, "provide http server port address")
	flag.StringVar(&config.AllowedOrigin, "allowed-origin", allowedOriginEnv, "provide allowed origin")
	flag.StringVar(&config.AccessTokenKey, "access-token-key", accessTokenKeyEnv, "provide access token secret key for jwt")
	flag.StringVar(&config.RefreshTokenKey, "refresh-token-key", refreshTokenKeyEnv, "provide refresh token secret key for jwt")
	flag.IntVar(&config.AccessTokenExpiration, "access-token-expiration", accessTokenExpirationEnv, "provide access token expiration time in seconds")
	flag.IntVar(&config.RefreshTokenExpiration, "refresh-token-expiration", refreshTokenExpirationEnv, "provide refresh token expiration time in seconds")
	flag.BoolVar(&config.AutoMigrate, "auto-migrate", autoMigrateEnv, "should auto migrate database (true | false)")

	flag.StringVar(&config.Postgres.Host, "postgres-host", postgresHost, "provide postgres host")
	flag.StringVar(&config.Postgres.Port, "postgres-port", postgresPort, "provide postgres port")
	flag.StringVar(&config.Postgres.DB, "postgres-db", postgresDB, "provide postgres db")
	flag.StringVar(&config.Postgres.User, "postgres-user", postgresUser, "provide postgres user")
	flag.StringVar(&config.Postgres.Password, "postgres-password", postgresPassword, "provide postgres password")
	flag.StringVar(&config.Postgres.SSLMode, "postgres-sslmode", postgresSSLModeEnv, "provide postgres ssl mode (disable | required)")

	flag.Parse()

	if os.Getenv("APP_ENV") == "dev" {
		strCfg, _ := json.MarshalIndent(&config, "", "  ")
		log.Println("Configuration:", string(strCfg))
	}
	return &config
}
