run:
	docker compose up -d
	sleep 5
	DB_USER=root \
	DB_PASS=password \
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_SCHEMA=bestmatch \
	HTTP_PORT=8080 \
	go run ./cmd/api/
