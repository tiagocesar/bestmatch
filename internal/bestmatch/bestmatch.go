package bestmatch

import (
	"context"

	"github.com/google/uuid"

	"github.com/tiagocesar/bestmatch/internal/models"
)

type matchRetriever interface {
	GetPartnerInfo(ctx context.Context, id uuid.UUID) (*models.Partner, error)
	GetMatches(ctx context.Context, req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error)
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

func (s *service) GetMatches(ctx context.Context,
	req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error) {
	// Validating the request model
	if err := req.Validate(); err != nil {
		return nil, err
	}

	partners, err := s.repo.GetMatches(ctx, req)
	if err != nil {
		return nil, err
	}

	return partners, nil
}
