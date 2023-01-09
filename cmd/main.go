package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/edwintantawi/taskit/cmd/config"
	authHTTPHandler "github.com/edwintantawi/taskit/internal/auth/delivery/http"
	authMiddleware "github.com/edwintantawi/taskit/internal/auth/delivery/http/middleware"
	authRepository "github.com/edwintantawi/taskit/internal/auth/repository"
	authUsecase "github.com/edwintantawi/taskit/internal/auth/usecase"
	taskHTTPHandler "github.com/edwintantawi/taskit/internal/task/delivery/http"
	taskRepository "github.com/edwintantawi/taskit/internal/task/repository"
	taskUsecase "github.com/edwintantawi/taskit/internal/task/usecase"
	userHTTPHandler "github.com/edwintantawi/taskit/internal/user/delivery/http"
	userRepository "github.com/edwintantawi/taskit/internal/user/repository"
	userUsecase "github.com/edwintantawi/taskit/internal/user/usecase"
	"github.com/edwintantawi/taskit/pkg/httpsvr"
	"github.com/edwintantawi/taskit/pkg/idgen"
	"github.com/edwintantawi/taskit/pkg/postgres"
	"github.com/edwintantawi/taskit/pkg/security"
	"github.com/edwintantawi/taskit/pkg/validator"
)

func main() {
	cfg := config.New()

	// Create new postgres connection.
	db, migrate := postgres.New(&cfg.Postgres)
	defer db.Close()

	// Migrate database.
	if err := migrate(cfg.AutoMigrate); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Create new providers.
	hashProvider := security.NewBcrypt()
	idProvider := idgen.NewUUID()
	validator := validator.New()
	jwtProvider := security.NewJWT(
		security.JWTTokenConfig{Key: cfg.AccessTokenKey, Exp: cfg.AccessTokenExpiration},
		security.JWTTokenConfig{Key: cfg.RefreshTokenKey, Exp: cfg.RefreshTokenExpiration},
	)

	// User.
	userRepository := userRepository.New(db, &idProvider)
	userUsecase := userUsecase.New(&validator, &userRepository, &hashProvider)
	userHTTPHandler := userHTTPHandler.New(&validator, &userUsecase)

	// Auth.
	authRepository := authRepository.New(db, &idProvider)
	authUsecase := authUsecase.New(&validator, &authRepository, &userRepository, &hashProvider, &jwtProvider)
	authHTTPHandler := authHTTPHandler.New(&validator, authUsecase)
	authMiddleware := authMiddleware.New(&jwtProvider)

	// Task.
	taskRepository := taskRepository.New(db, &idProvider)
	taskUsecase := taskUsecase.New(taskRepository)
	taskHTTPHandler := taskHTTPHandler.New(&validator, taskUsecase)

	// Create new router.
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.AllowedOrigin},
		AllowedMethods:   []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/users", userHTTPHandler.Post)

		r.Post("/api/authentications", authHTTPHandler.Post)
		r.Put("/api/authentications", authHTTPHandler.Put)
	})

	// private routes (need authentication)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Get("/api/authentications", authHTTPHandler.Get)
		r.Delete("/api/authentications", authHTTPHandler.Delete)

		r.Post("/api/tasks", taskHTTPHandler.Post)
		r.Get("/api/tasks", taskHTTPHandler.Get)
		r.Get("/api/tasks/{task_id}", taskHTTPHandler.GetByID)
		r.Delete("/api/tasks/{task_id}", taskHTTPHandler.Delete)
		r.Put("/api/tasks/{task_id}", taskHTTPHandler.Put)
	})

	// Start HTTP server.
	log.Printf("Server running at %s", cfg.Port)
	svr := httpsvr.New(":"+cfg.Port, r)
	if err := svr.Run(); err != nil {
		log.Fatal(err)
	}
}
