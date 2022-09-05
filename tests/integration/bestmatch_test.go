// //go:build integration

package integration

import (
	"context"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/tiagocesar/bestmatch/internal/bestmatch"
	"github.com/tiagocesar/bestmatch/internal/models"
	"github.com/tiagocesar/bestmatch/internal/repo"
)

func Test_GetMatches_Validations(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "success"},
		{name: "error - no materials specified"},
		{name: "error - no address provided"},
		{name: "error - no phone number"},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {

		})
	}
}

func Test_GetMatches_Bestmatch(t *testing.T) {
	r, err := repo.NewRepository("root", "password", "localhost", "5432", "bestmatch")
	if err != nil {
		log.Fatal(err)
	}

	svc := bestmatch.NewService(r)

	t.Run("Two matches, ordered by best match", func(t *testing.T) {
		req := getMatchRequestWithDefaults()

		req.Materials = append(req.Materials,
			"07cab731-d981-4915-9444-cc997eec351f",
			"1606f175-3502-4028-9501-6b591c00f1f3",
			"ac47d822-ffc9-48b7-8492-4d49e921d4df")

		req.Address = struct {
			Lat  string `json:"lat"`
			Long string `json:"long"`
		}(struct {
			Lat  string
			Long string
		}{Lat: "52.3599795", Long: "4.8851198"})

		m, err := svc.GetMatches(context.Background(), req)
		require.NoError(t, err)

		matches := *m
		require.Equal(t, 2, len(matches))
		require.Equal(t, "b276cb54-ac52-4f8c-adb1-afce5ced67c4", matches[0].Id.String())
	})
}

func getMatchRequestWithDefaults() models.ListPartnersByMatchRequest {
	return models.ListPartnersByMatchRequest{
		Area:        100.0,
		PhoneNumber: "0123456789",
	}
}
