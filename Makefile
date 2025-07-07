test:: pg_up go_test  pg_down

run:: pg_up go_run

go_test:: go_tidy statictest
	go test ./...

go_run::
	go run cmd/gophermart/main.go \
	 -d "postgres://gofermart:yandex@localhost:5432/gofermart_db?sslmode=disable"

statictest::
	go vet -vettool=$(shell which statictest) ./...


pg_up::
	docker compose up -d

pg_down::
	docker compose down -v

go_tidy::
	go mod tidy


prepare::
	go install github.com/pressly/goose/v3/cmd/goose@latest

migration_create:: prepare
	goose -dir internal/repo/migrator/migrations create migration sql
