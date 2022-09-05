package models

import (
	"github.com/google/uuid"
)

type Partner struct {
	Id      uuid.UUID `db,json:"id"`
	Name    string    `db,json:"name"`
	Address string    `db,json:"address"`
	Radius  int       `db,json:"radius"`
	Rating  int       `db,json:"rating"`
}

// PartnerResult enriches Partner info with extra query-specific information
type PartnerResult struct {
	Ranking  int     `json:"ranking"`
	Distance float32 `json:"distance_km"` // Distance in km
	Partner
}
