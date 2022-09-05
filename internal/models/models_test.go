package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_Validate tests Validate() for the ListPartnersByMatchRequest struct
func Test_Validate(t *testing.T) {
	defaultReq := func() ListPartnersByMatchRequest {
		return ListPartnersByMatchRequest{
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

	tests := []struct {
		name        string
		req         func() ListPartnersByMatchRequest
		expectedErr error
	}{
		{
			name: "success",
			req:  defaultReq,
		},
		{
			name: "invalid materials",
			req: func() ListPartnersByMatchRequest {
				r := defaultReq()
				r.Materials = []string{}

				return r
			},
			expectedErr: errors.New("no materials provided"),
		},
		{
			name: "invalid area",
			req: func() ListPartnersByMatchRequest {
				r := defaultReq()
				r.Area = 0

				return r
			},
			expectedErr: errors.New("no area provided"),
		},
		{
			name: "invalid phone number",
			req: func() ListPartnersByMatchRequest {
				r := defaultReq()
				r.PhoneNumber = ""

				return r
			},
			expectedErr: errors.New("no phone number provided"),
		},
		{
			name: "invalid address - latitude",
			req: func() ListPartnersByMatchRequest {
				r := defaultReq()
				r.Address.Lat = ""

				return r
			},
			expectedErr: errors.New("no address provided"),
		},
		{
			name: "invalid address - longitude",
			req: func() ListPartnersByMatchRequest {
				r := defaultReq()
				r.Address.Long = ""

				return r
			},
			expectedErr: errors.New("no address provided"),
		},
		{
			name: "invalid address - missing latitude and longitude",
			req: func() ListPartnersByMatchRequest {
				r := defaultReq()
				r.Address.Lat = ""
				r.Address.Long = ""

				return r
			},
			expectedErr: errors.New("no address provided"),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.req().Validate()

			require.Equal(t, test.expectedErr, err)
		})
	}
}
