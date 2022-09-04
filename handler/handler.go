package handler

import (
	"context"

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
