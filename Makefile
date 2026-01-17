# Default to local settings
DB_URL ?= postgresql://admin:password@localhost:5432/app_db?sslmode=disable

migrate-up:
	migrate -database "$(DB_URL)" -path db/migrations up

migrate-down:
	migrate -database "$(DB_URL)" -path db/migrations down

# Helpful for staging/production (Usage: make migrate-remote URL=...)
migrate-remote:
	migrate -database "$(URL)" -path db/migrations up

jet-generate:
	jet -source=PostgreSQL -host=localhost -port=5432 -user=admin -password=password -dbname=app_db -path=./db/model -schema=public
