create_migration:
# make create_migration name=name_of_migration_without_spaces
	migrate create -ext sql -dir db/migrations -seq ${name}
migrate:
	migrate -database 'postgres://postgres:1234@localhost:5432/postgres?sslmode=disable' -path ./db/migrations/ up
db_win:
	docker run -d --name pgsql --hostname db -p 5432:5432 -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -v C:\ProgramFiles\PostgreSQL\14\data:/var/lib/postgresql/data --network=net_postgres postgres
build_image:
	docker build -t i_samplicity .
local:
	go run cmd/main.go