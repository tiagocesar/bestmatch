package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/tiagocesar/bestmatch/internal/models"
)

type matchRetriever interface {
	GetPartnerInfo(ctx context.Context, id string) (*models.Partner, error)
}

type handler struct {
	service matchRetriever
}

func NewHandler(service matchRetriever) *handler {
	return &handler{service: service}
}

func (h *handler) ConfigureAndServe(port string) {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	router.Get("/health", health)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatalf("Failed to start HTTP server")
	}
}

func health(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "ok")
	w.WriteHeader(http.StatusOK)
}
