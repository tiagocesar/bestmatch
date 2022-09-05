package bestmatch

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tiagocesar/bestmatch/internal/models"
)

type repoFake struct {
	GetPartnerInfoFn func(id uuid.UUID) (*models.Partner, error)
	GetMatchesFn     func(req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error)
}

func (r repoFake) GetPartnerInfo(_ context.Context, id uuid.UUID) (*models.Partner, error) {
	return r.GetPartnerInfoFn(id)
}

func (r repoFake) GetMatches(ctx context.Context,
	req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error) {
	return r.GetMatchesFn(req)
}

func Test_GetPartnerInfo(t *testing.T) {
	u := uuid.New()

	tests := []struct {
		name string
		id   string
		repo repoFake
	}{
		{
			name: "success",
			id:   u.String(),
			repo: repoFake{GetPartnerInfoFn: func(id uuid.UUID) (*models.Partner, error) {
				return &models.Partner{}, nil
			}},
		},
		{
			name: "invalid uuid",
			id:   "invalid string",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			svc := NewService(test.repo)

			_, err := svc.GetPartnerInfo(context.Background(), test.id)

			if err != nil {
				require.True(t, uuid.IsInvalidLengthError(err))
				return
			}

			assert.NoError(t, err)
		})
	}
}

func Test_GetMatches(t *testing.T) {
	defaultRequest := models.ListPartnersByMatchRequest{
		Materials:   []string{"asd"},
		Area:        100.0,
		PhoneNumber: "0123456789",
		Address: struct {
			Lat  string `json:"lat"`
			Long string `json:"long"`
		}(struct {
			Lat  string
			Long string
		}{Lat: "0.0", Long: "0.0"}),
	}

	errRequest := defaultRequest
	errRequest.Area = 0

	tests := []struct {
		name          string
		request       models.ListPartnersByMatchRequest
		repo          repoFake
		expectedError error
	}{
		{
			name:    "success",
			request: defaultRequest,
			repo: repoFake{GetMatchesFn: func(req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error) {
				return &[]models.PartnerResult{}, nil
			}},
		},
		{
			name:          "request validation error",
			request:       errRequest,
			expectedError: errors.New("no area provided"),
		},
		{
			name: "repository error", request: defaultRequest,
			repo: repoFake{GetMatchesFn: func(req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error) {
				return nil, sql.ErrNoRows
			}},
			expectedError: sql.ErrNoRows,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			svc := service{repo: test.repo}

			_, err := svc.GetMatches(context.Background(), test.request)

			require.Equal(t, test.expectedError, err)
		})
	}
}
