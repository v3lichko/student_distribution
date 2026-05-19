-include Makefile.mk
export

APP_PORT       ?= 8080
POSTGRES_HOST  ?= localhost
POSTGRES_PORT  ?= 5433
POSTGRES_USER  ?=
POSTGRES_PASSWORD ?=
POSTGRES_DB    ?=
CSV_PATH       ?= raw_data/students.csv
EXPORT_PATH    ?= raw_data/distribution_export.csv
GROUP_COUNT    ?= 10
GROUP_CAPACITY ?= 1000

.PHONY: init run db-up db-down migrate migrate-students migrate-groups migrate-assign \
        generate-csv import-csv run-distribution export-distribution swagger help

help:
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | awk -F':.*?## ' '{printf "  %-24s %s\n", $$1, $$2}'

init:
	@test -f Makefile.mk || cp Makefile.mk.dist Makefile.mk
	@test -f .env        || cp .env.dist .env
	@echo "init: Makefile.mk and .env are ready — fill in credentials"

run:
	go run ./cmd/api/main.go

db-up:
	docker compose up -d

db-down:
	docker compose down

PSQL = PGPASSWORD=$(POSTGRES_PASSWORD) psql -h $(POSTGRES_HOST) -p $(POSTGRES_PORT) -U $(POSTGRES_USER) -d $(POSTGRES_DB)

migrate: migrate-students migrate-groups migrate-assign

migrate-students:
	$(PSQL) -f migrations/create_students.up.sql

migrate-groups:
	$(PSQL) -f migrations/create_groups.up.sql

migrate-assign:
	$(PSQL) -f migrations/students_to_group.up.sql

generate-csv:        
	python3 scripts/generate_csv.py

import-csv:          
	curl -X POST http://localhost:$(APP_PORT)/students/import -F "file=@$(CSV_PATH)"

run-distribution:
	curl -X POST http://localhost:$(APP_PORT)/distribution/run

export-distribution:
	curl http://localhost:$(APP_PORT)/distribution/export -o $(EXPORT_PATH)

swagger:
	swag init -g cmd/api/main.go
