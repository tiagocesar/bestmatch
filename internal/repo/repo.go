package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

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

// GetPartnerInfo gets info about a specific partner, identified by its id
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

// GetMatches will get an ordered list of partners, ordered by "best match".
//
// Best match is determined by the average rating of the partner, and then by distance to the customer.
// To be considered a match, a partner needs to be experienced with all materials present in the request.
//
// Matches that are out of the radius of operation of the company shouldn't be considered.
func (r *repository) GetMatches(ctx context.Context,
	req models.ListPartnersByMatchRequest) (*[]models.PartnerResult, error) {

	// uuids need to be between quotes, so we need to add them to each element in the
	// req.Materials array, then convert all elements to a single string that will
	// be concatenated against the query
	for i := range req.Materials {
		req.Materials[i] = "'" + req.Materials[i] + "'"
	}
	materials := strings.Join(req.Materials, ",")

	q := `
		SELECT p.id, p.name, p.address, p.radius, p.rating, p1.distance_km
		  FROM (SELECT p.id, (point($1,$2) <@> p.address) * 1.609344 as distance_km
				  FROM partners p
				  JOIN partners_materials pm on p.id = pm.partner_id
				 WHERE pm.material_id IN (:materials:)
				 GROUP BY p.id
				HAVING COUNT(pm.material_id) = $3) p1
		  JOIN partners p ON p.id = p1.id
		 WHERE distance_km <= p.radius
		 ORDER BY p.rating DESC, distance_km DESC`
	q = strings.ReplaceAll(q, ":materials:", materials)

	rows, err := r.db.QueryContext(ctx, q, req.Address.Long, req.Address.Lat, len(req.Materials))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) { _ = rows.Close() }(rows)

	results := make([]models.PartnerResult, 0)
	ranking := 1
	for rows.Next() {
		var pr models.PartnerResult
		err := rows.Scan(&pr.Id, &pr.Name, &pr.Address, &pr.Radius, &pr.Rating, &pr.Distance)
		if err != nil {
			return nil, err
		}
		pr.Ranking = ranking

		results = append(results, pr)
		ranking++
	}

	return &results, nil
}
