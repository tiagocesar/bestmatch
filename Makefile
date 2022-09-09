run:
	docker compose up -d --wait
	DB_USER=root \
	DB_PASS=password \
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_SCHEMA=bestmatch \
	HTTP_PORT=8080 \
	go run ./cmd/api/

integration-tests:
	docker compose down --volumes
	sleep 2
	docker compose up -d --wait
	go test -v -tags=integration ./...
