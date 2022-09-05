package models

import (
	"errors"

	"github.com/google/uuid"
)

type Partner struct {
	Id      uuid.UUID `db,json:"id"`
	Name    string    `db,json:"name"`
	Address string    `db,json:"address"`
	Radius  int       `db,json:"radius"`
	Rating  float32   `db,json:"rating"`
}

// PartnerResult enriches Partner info with extra query-specific information
type PartnerResult struct {
	Ranking  int     `json:"ranking"`
	Distance float32 `json:"distance_km"` // Distance in km
	Partner
}

// ListPartnersByMatchRequest represents a customer request that will be used to find the best matches
type ListPartnersByMatchRequest struct {
	Materials   []string `json:"materials"` // The ID's of the desired materials
	Area        float32  `json:"area"`      // Construction area in square meters
	PhoneNumber string   `json:"phone"`
	Address     struct {
		Lat  string `json:"lat"`
		Long string `json:"long"`
	} `json:"address"`
}

func (l ListPartnersByMatchRequest) Validate() error {
	if len(l.Materials) == 0 {
		return errors.New("no materials provided")
	}

	if l.Area == 0 {
		return errors.New("no area provided")
	}

	if l.PhoneNumber == "" {
		return errors.New("no phone number provided")
	}

	if l.Address.Lat == "" || l.Address.Long == "" {
		return errors.New("no address provided")
	}

	return nil
}
