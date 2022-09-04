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
