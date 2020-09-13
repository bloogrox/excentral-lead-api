.PHONY: psql
psql:
	docker-compose run --rm postgres psql --host=postgres --dbname=postgres --username=postgres


.PHONY: gogenerate
gogenerate:
	docker-compose run --rm test go generate ./...


.PHONY: test
test:
	docker-compose run --rm test go test ./...


.PHONY: migrate
migrate:
	docker-compose run --rm web go run cmd/migrate/main.go
