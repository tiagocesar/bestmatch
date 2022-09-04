package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/tiagocesar/bestmatch/handler"
	"github.com/tiagocesar/bestmatch/internal/bestmatch"
	"github.com/tiagocesar/bestmatch/internal/repo"
)

const (
	EnvDbUser   = "DB_USER"
	EnvDbPass   = "DB_PASS"
	EnvDbHost   = "DB_HOST"
	EnvDbPort   = "DB_PORT"
	EnvDbSchema = "DB_SCHEMA"

	EnvHttpPort = "HTTP_PORT"
)

func main() {
	// Getting environment vars
	envVars, err := getEnvVars()
	if err != nil {
		log.Fatal(err)
	}

	r, err := repo.NewRepository(envVars[EnvDbUser], envVars[EnvDbPass], envVars[EnvDbHost],
		envVars[EnvDbPort], envVars[EnvDbSchema])
	if err != nil {
		log.Fatal(err)
	}

	svc := bestmatch.NewService(r)

	h := handler.NewHandler(svc)

	log.Printf("HTTP server starting on port %s", envVars[EnvHttpPort])
	h.ConfigureAndServe(envVars[EnvHttpPort])

	log.Println("HTTP server exiting")
	os.Exit(0)
}

// getEnvVars gets all environment variables necessary for this service to run.
func getEnvVars() (map[string]string, error) {
	result := make(map[string]string)
	var ok bool

	// DB vars
	if result[EnvDbUser], ok = os.LookupEnv(EnvDbUser); !ok {
		return nil, errors.New(fmt.Sprintf("environment variable %s not set", EnvDbUser))
	}

	if result[EnvDbPass], ok = os.LookupEnv(EnvDbPass); !ok {
		return nil, errors.New(fmt.Sprintf("environment variable %s not set", EnvDbPass))
	}

	if result[EnvDbHost], ok = os.LookupEnv(EnvDbHost); !ok {
		return nil, errors.New(fmt.Sprintf("environment variable %s not set", EnvDbHost))
	}

	if result[EnvDbPort], ok = os.LookupEnv(EnvDbPort); !ok {
		return nil, errors.New(fmt.Sprintf("environment variable %s not set", EnvDbPort))
	}

	if result[EnvDbSchema], ok = os.LookupEnv(EnvDbSchema); !ok {
		return nil, errors.New(fmt.Sprintf("environment variable %s not set", EnvDbSchema))
	}

	if result[EnvHttpPort], ok = os.LookupEnv(EnvHttpPort); !ok {
		return nil, errors.New(fmt.Sprintf("environment variable %s not set", EnvHttpPort))
	}

	return result, nil
}
