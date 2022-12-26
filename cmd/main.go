package main

import (
	"log"

	"github.com/go-chi/chi/v5"

	"github.com/edwintantawi/taskit/cmd/config"
	"github.com/edwintantawi/taskit/internal/user"
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

	userRepository := user.NewRepository(db, idProvider)
	userUsecase := user.NewUsecase(userRepository, hashProvider)
	userHTTPHandler := user.NewHTTPHandler(userUsecase)

	r := chi.NewRouter()
	r.Post("/api/users", userHTTPHandler.Post)

	log.Printf("Server running at %s", cfg.ServerAddr)
	svr := httpsvr.New(cfg.ServerAddr, r)
	if err := svr.Run(); err != nil {
		log.Fatal(err)
	}
}
