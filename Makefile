migrate-up:
	migrate -database "postgresql://admin:password@localhost:5432/app_db?sslmode=disable" -path db/migrations up

migrate-down:
	migrate -database "postgresql://admin:password@localhost:5432/app_db?sslmode=disable" -path db/migrations down
