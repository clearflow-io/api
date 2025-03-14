migrate-up:
	migrate -database "postgresql://admin:password@localhost:5432/app_db?sslmode=disable" -path db/migrations up

migrate-down:
	migrate -database "postgresql://admin:password@localhost:5432/app_db?sslmode=disable" -path db/migrations down

jet-generate:
	jet -source=PostgreSQL -host=localhost -port=5432 -user=admin -password=password -dbname=app_db -path=./db/model -schema=public
