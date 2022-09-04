package bestmatch

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tiagocesar/bestmatch/internal/models"
)

type repoFake struct {
	GetPartnerInfoFn func(id uuid.UUID) (*models.Partner, error)
}

func (r repoFake) GetPartnerInfo(id uuid.UUID) (*models.Partner, error) {
	return r.GetPartnerInfoFn(id)
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

			_, err := svc.GetPartnerInfo(test.id)

			if err != nil {
				require.True(t, uuid.IsInvalidLengthError(err))
				return
			}

			assert.NoError(t, err)
		})
	}
}
