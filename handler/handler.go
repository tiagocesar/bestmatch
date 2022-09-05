package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/tiagocesar/bestmatch/internal/models"
)

type matchRetriever interface {
	GetPartnerInfo(ctx context.Context, id string) (*models.Partner, error)
	GetMatches(ctx context.Context, req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error)
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
	router.Use(middleware.StripSlashes)

	router.Get("/health", health)
	router.Get("/partners", h.listPartnersByMatch)
	router.Get("/partners/{id}", h.getPartner)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		log.Fatalf("Failed to start HTTP server")
	}
}

func health(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "ok")
	w.WriteHeader(http.StatusOK)
}

func (h *handler) listPartnersByMatch(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var matchRequest models.ListPartnersByMatchRequest
	err := json.NewDecoder(req.Body).Decode(&matchRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	partners, err := h.service.GetMatches(ctx, matchRequest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, _ := json.Marshal(partners)

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(j))
}

// getPartner gets information about a specific partner
func (h *handler) getPartner(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := chi.URLParam(req, "id")

	partner, err := h.service.GetPartnerInfo(ctx, id)
	switch {
	case err == nil:
		break
	case uuid.IsInvalidLengthError(err):
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid id"))
		return
	case errors.Is(err, sql.ErrNoRows):
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, _ := json.Marshal(partner)

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, string(j))
}
