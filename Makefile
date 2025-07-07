test:: pg_up go_test  pg_down


go_test:: statictest
	go test ./...

statictest::
	go vet -vettool=$(shell which statictest) ./...


pg_up::
	docker compose up -d

pg_down::
	docker compose down -v
