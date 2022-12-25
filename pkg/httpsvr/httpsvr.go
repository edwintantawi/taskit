package httpsvr

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Addr   string
	Router *chi.Mux
}

func New(addr string, r *chi.Mux) *Server {
	return &Server{Addr: addr, Router: r}
}

// Run starts the HTTP server with gracefully shutdown.
func (s *Server) Run() error {
	svr := &http.Server{
		Addr:    s.Addr,
		Handler: s.Router,
	}

	shutdownChan := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("Shutting down the server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownChan <- svr.Shutdown(ctx)
	}()

	err := svr.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return <-shutdownChan
}
