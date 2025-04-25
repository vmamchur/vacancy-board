build:
	docker compose build 

run:	build
	docker compose up

migrate-up:
	goose -dir db/migrations postgres postgres://postgres:12345@localhost:5432/postgres?sslmode=disable up

migrate-down:
	goose -dir db/migrations postgres postgres://postgres:12345@localhost:5432/postgres?sslmode=disable down

