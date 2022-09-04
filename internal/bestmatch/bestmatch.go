package bestmatch

import (
	"context"

	"github.com/google/uuid"

	"github.com/tiagocesar/bestmatch/internal/models"
)

type matchRetriever interface {
	GetPartnerInfo(ctx context.Context, id uuid.UUID) (*models.Partner, error)
}

type service struct {
	repo matchRetriever
}

func NewService(repo matchRetriever) *service {
	return &service{repo: repo}
}

func (s *service) GetPartnerInfo(ctx context.Context, id string) (*models.Partner, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	p, err := s.repo.GetPartnerInfo(ctx, u)
	if err != nil {
		return nil, err
	}

	return p, nil
}
