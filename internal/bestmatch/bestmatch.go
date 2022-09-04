package bestmatch

import (
	"github.com/google/uuid"

	"github.com/tiagocesar/bestmatch/internal/models"
)

type matchRetriever interface {
	GetPartnerInfo(id uuid.UUID) (*models.Partner, error)
}

type service struct {
	repo matchRetriever
}

func NewService(repo matchRetriever) *service {
	return &service{repo: repo}
}

func (s *service) GetPartnerInfo(id string) (*models.Partner, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	p, err := s.repo.GetPartnerInfo(u)
	if err != nil {
		return nil, err
	}

	return p, nil
}
