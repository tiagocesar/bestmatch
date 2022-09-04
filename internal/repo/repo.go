package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"github.com/tiagocesar/bestmatch/internal/models"
)

type repository struct {
	db *sql.DB
}

func NewRepository(user, pass, host, port, schema string) (*repository, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, schema)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &repository{db: db}, nil
}

func (r *repository) GetPartnerInfo(ctx context.Context, id uuid.UUID) (*models.Partner, error) {
	q := `SELECT id, name, address, radius, rating
            FROM partners
           WHERE id = $1`

	var partner models.Partner
	err := r.db.QueryRowContext(ctx, q, id).Scan(&partner.Id, &partner.Name, &partner.Address, &partner.Radius,
		&partner.Rating)
	if err != nil {
		return nil, err
	}

	return &partner, nil
}
