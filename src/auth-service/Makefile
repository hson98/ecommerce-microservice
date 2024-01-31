local:
	echo "Starting local environment"
	docker-compose -f docker-compose.local.yml up --build

force:
	migrate -database postgres://postgres:postgres@localhost:5435/auth_db?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://postgres:postgres@localhost:5435/auth_db?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5435/auth_db?sslmode=disable -path migrations up

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5435/auth_db?sslmode=disable -path migrations down 1