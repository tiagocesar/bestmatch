package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/tiagocesar/bestmatch/internal/models"
)

type matchRetrieverFake struct {
	GetPartnerInfoFn func(id string) (*models.Partner, error)
}

func (m matchRetrieverFake) GetPartnerInfo(ctx context.Context, id string) (*models.Partner, error) {
	return m.GetPartnerInfoFn(id)
}

func Test_getPartner(t *testing.T) {
	tests := []struct {
		name             string
		expectedRespCode int
		matchRetriever   matchRetrieverFake
		expectedResponse string
	}{
		{
			name:             "success",
			expectedRespCode: http.StatusOK,
			matchRetriever: matchRetrieverFake{GetPartnerInfoFn: func(id string) (*models.Partner, error) {
				return &models.Partner{}, nil
			}},
		},
		{
			name:             "invalid id format",
			expectedRespCode: http.StatusBadRequest,
			matchRetriever: matchRetrieverFake{func(id string) (*models.Partner, error) {
				_, err := uuid.Parse("invalid string")
				return nil, err
			}},
			expectedResponse: "invalid id",
		},
		{
			name:             "partner not found",
			expectedRespCode: http.StatusNotFound,
			matchRetriever: matchRetrieverFake{func(id string) (*models.Partner, error) {
				return nil, sql.ErrNoRows
			}},
		},
		{
			name:             "internal server error",
			expectedRespCode: http.StatusInternalServerError,
			matchRetriever: matchRetrieverFake{func(id string) (*models.Partner, error) {
				return nil, errors.New("random internal error")
			}},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("id", "")

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/partners/{id}", nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
			require.NoError(t, err)

			h := handler{service: test.matchRetriever}

			h.getPartner(rr, req)

			require.Equal(t, test.expectedRespCode, rr.Code)

			body := rr.Body.String()

			if body != "" {
				expected := test.expectedResponse

				if expected == "" {
					j, _ := json.Marshal(&models.Partner{})
					expected = string(j)
				}

				require.Equal(t, expected, body)
			}
		})
	}
}
