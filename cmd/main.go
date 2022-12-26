package main

import (
	"log"

	"github.com/go-chi/chi/v5"

	"github.com/edwintantawi/taskit/cmd/config"
	authHTTPHandler "github.com/edwintantawi/taskit/internal/auth/delivery/http"
	authRepository "github.com/edwintantawi/taskit/internal/auth/repository"
	authUsecase "github.com/edwintantawi/taskit/internal/auth/usecase"
	userHTTPHandler "github.com/edwintantawi/taskit/internal/user/delivery/http"
	userRepository "github.com/edwintantawi/taskit/internal/user/repository"
	userUsecase "github.com/edwintantawi/taskit/internal/user/usecase"
	"github.com/edwintantawi/taskit/pkg/httpsvr"
	"github.com/edwintantawi/taskit/pkg/idgen"
	"github.com/edwintantawi/taskit/pkg/postgres"
	"github.com/edwintantawi/taskit/pkg/security"
)

func main() {
	cfg, err := config.New(".")
	if err != nil {
		log.Fatal(err)
	}

	db := postgres.New(cfg.PostgresDSN)
	defer db.Close()

	hashProvider := security.NewBcrypt()
	idProvider := idgen.NewUUID()
	jwtProvider := security.NewJWT(
		security.JWTTokenConfig{Key: cfg.AccessTokenKey, Exp: cfg.AccessTokenExp},
		security.JWTTokenConfig{Key: cfg.RefreshTokenKey, Exp: cfg.RefreshTokenExp},
	)

	userRepository := userRepository.New(db, idProvider)
	userUsecase := userUsecase.New(userRepository, hashProvider)
	userHTTPHandler := userHTTPHandler.New(userUsecase)

	authRepository := authRepository.New(db, idProvider)
	authUsecase := authUsecase.New(authRepository, userRepository, hashProvider, jwtProvider)
	authHTTPHandler := authHTTPHandler.New(authUsecase)

	r := chi.NewRouter()
	r.Post("/api/users", userHTTPHandler.Post)
	r.Post("/api/authentications", authHTTPHandler.Post)

	log.Printf("Server running at %s", cfg.ServerAddr)
	svr := httpsvr.New(cfg.ServerAddr, r)
	if err := svr.Run(); err != nil {
		log.Fatal(err)
	}
}
