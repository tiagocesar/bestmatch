package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/tiagocesar/bestmatch/internal/models"
)

type matchRetrieverFake struct {
	GetPartnerInfoFn func(id string) (*models.Partner, error)
	GetMatchesFn     func(req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error)
}

func (m matchRetrieverFake) GetPartnerInfo(_ context.Context, id string) (*models.Partner, error) {
	return m.GetPartnerInfoFn(id)
}

func (m matchRetrieverFake) GetMatches(_ context.Context,
	req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error) {
	return m.GetMatchesFn(req)
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
			matchRetriever: matchRetrieverFake{GetPartnerInfoFn: func(id string) (*models.Partner, error) {
				_, err := uuid.Parse("invalid string")
				return nil, err
			}},
			expectedResponse: "invalid id",
		},
		{
			name:             "partner not found",
			expectedRespCode: http.StatusNotFound,
			matchRetriever: matchRetrieverFake{GetPartnerInfoFn: func(id string) (*models.Partner, error) {
				return nil, sql.ErrNoRows
			}},
		},
		{
			name:             "internal server error",
			expectedRespCode: http.StatusInternalServerError,
			matchRetriever: matchRetrieverFake{GetPartnerInfoFn: func(id string) (*models.Partner, error) {
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

func Test_listPartnersByMatch(t *testing.T) {
	j, _ := json.Marshal(getMatchRequest())
	getMatchRequestStr := string(j)

	j, _ = json.Marshal(&[]models.PartnerResult{{}})
	partnerResultStr := string(j)

	tests := []struct {
		name           string
		requestBody    string
		matchRetriever matchRetrieverFake

		expectedRespCode int    // Always checked
		expectedResponse string // Not checked if empty
	}{
		{
			name:        "success",
			requestBody: getMatchRequestStr,
			matchRetriever: matchRetrieverFake{GetMatchesFn: func(req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error) {
				return &[]models.PartnerResult{
					{},
				}, nil
			}},
			expectedRespCode: http.StatusOK,
			expectedResponse: partnerResultStr,
		},
		{
			name:             "invalid request body",
			expectedRespCode: http.StatusInternalServerError,
		},
		{
			name:        "internal server error from the repo",
			requestBody: getMatchRequestStr,
			matchRetriever: matchRetrieverFake{GetMatchesFn: func(req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error) {
				return nil, sql.ErrNoRows
			}},
			expectedRespCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, "/partners", strings.NewReader(test.requestBody))
			require.NoError(t, err)

			h := handler{service: test.matchRetriever}

			h.listPartnersByMatch(rr, req)

			require.Equal(t, test.expectedRespCode, rr.Code)

			body := rr.Body.String()

			if body != "" {
				expected := test.expectedResponse

				if expected == "" {
					j, _ := json.Marshal(&models.PartnerResult{})
					expected = string(j)
				}

				require.Equal(t, expected, body)
			}
		})
	}
}

func getMatchRequest() models.ListPartnersByMatchRequest {
	return models.ListPartnersByMatchRequest{
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
}
