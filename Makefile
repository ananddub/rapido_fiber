
include .env

export PGPASSWORD := $(POSTGRES_PASSWORD)

pgcli:
	@pgcli  -U $(POSTGRES_USER) -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) $(POSTGRES_DB)

mz:
	@psql -h $(MATERIALIZE_HOST) -U $(MATERIALIZE_USER) -p $(MATERIALIZE_PORT)

redis:
	@iredis -h $(REDIS_HOST) -p $(REDIS_PORT)

rd:
	@docker exec -it rapido-redpanda bash

run:
	@go run cmd/main.go
