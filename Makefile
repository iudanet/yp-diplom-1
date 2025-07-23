test:: mock_gen lint pg_up go_test  test_ci pg_down

run:: lint pg_up go_run

run_accrual::
	    ./cmd/accrual/accrual_linux_amd64 -a :8081 -d "postgresql://gofermart:yandex@127.0.0.1/gofermart_db?sslmode=disable"

lint:: statictest golangci-fmt  golangci

go_test:: go_tidy statictest
	go test -race ./...

go_run::
	go run cmd/gophermart/main.go \
	 -d "postgres://gofermart:yandex@localhost:5432/gofermart_db?sslmode=disable"

statictest::
	go vet -vettool=$(shell which statictest) ./...

mock_gen::
	mockgen -source=internal/service/models.go -destination=internal/service/mock_service/mock_service.go -package=mock_service
	mockgen -source=internal/repo/repo.go -destination=internal/repo/mock_repo/mock_repo.go -package=mock_repo
	mockgen -source=internal/service/accrual.go -destination=internal/service/mock_service/accrual_client_mock.go -package=mock_service

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

build::
	go build -o cmd/gophermart/gophermart cmd/gophermart/main.go

test_ci:: pg_up build
	gophermarttest \
	  -test.v -test.run=^TestGophermart$ \
	  -gophermart-binary-path=cmd/gophermart/gophermart \
	  -gophermart-host=localhost \
	  -gophermart-port=8080 \
	  -gophermart-database-uri="postgresql://gofermart:yandex@127.0.0.1/gofermart_db?sslmode=disable" \
	  -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
	  -accrual-host=localhost \
	  -accrual-port=$(shell random unused-port) \
	  -accrual-database-uri="postgresql://gofermart:yandex@127.0.0.1/gofermart_db?sslmode=disable"

test_ci_one:: pg_up build
	gophermarttest \
	  -test.v -test.run=^TestGophermart/TestUserOrders/order_upload$ \
	  -gophermart-binary-path=cmd/gophermart/gophermart \
	  -gophermart-host=localhost \
	  -gophermart-port=8080 \
	  -gophermart-database-uri="postgresql://gofermart:yandex@127.0.0.1/gofermart_db?sslmode=disable" \
	  -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
	  -accrual-host=localhost \
	  -accrual-port=$(shell random unused-port) \
	  -accrual-database-uri="postgresql://gofermart:yandex@127.0.0.1/gofermart_db?sslmode=disable"
