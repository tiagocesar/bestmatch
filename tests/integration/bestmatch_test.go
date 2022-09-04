// //go:build integration

package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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
	tests := []struct {
		name string
	}{
		{name: "no matches found"},
		{name: "one match in range"},
		{name: "two matches, ordered by rating"},
	}

	// Establishing a connection against the (containerized) test db
	dbContainer, db := setupTestDb()
	defer dbContainer.Terminate(context.Background())

	_ = db

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {

		})
	}
}

// setupTestDb spins up a container running postgres, with some data loaded, for integration testing purposes
// it's used so there's no need to define a docker-compose file for testing purposes
func setupTestDb() (testcontainers.Container, *sql.DB) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "bestmatch",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	err := dbContainer.CopyFileToContainer(context.Background(), "./scripts/db", "/docker-entrypoint-initdb.d/", 700)
	if err != nil {
		panic("Failed to copy init db file to container")
	}

	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	dbURI := fmt.Sprintf("postgres://postgres:postgres@%v:%v/bestmatch?sslmode=disable", host, port.Port())
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return dbContainer, db
}
