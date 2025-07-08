test:: lint pg_up go_test  pg_down

run:: lint pg_up go_run

lint:: statictest golangci-fmt  golangci

go_test:: go_tidy statictest
	go test ./...

go_run::
	go run cmd/gophermart/main.go \
	 -d "postgres://gofermart:yandex@localhost:5432/gofermart_db?sslmode=disable"

statictest::
	go vet -vettool=$(shell which statictest) ./...

golangci::
	golangci-lint run ./...

golangci-fmt::
	golangci-lint fmt ./...

pg_up::
	docker compose up -d

pg_down::
	docker compose down -v

go_tidy::
	go mod tidy


prepare::
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.1

migration_create:: prepare
	goose -dir internal/repo/migrator/migrations create migration sql
