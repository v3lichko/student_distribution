-include Makefile.mk
export

APP_PORT     ?= 8080
POSTGRES_HOST ?= localhost
POSTGRES_PORT ?= 5433
POSTGRES_USER ?= jigoku
POSTGRES_PASSWORD ?= 123
POSTGRES_DB  ?= students_db
CSV_PATH       ?= raw_data/students.csv
EXPORT_PATH    ?= raw_data/distribution_export.csv
GROUP_COUNT    ?= 10
GROUP_CAPACITY ?= 1000

.PHONY: init run db-up db-down migrate migrate-students migrate-groups migrate-assign \
		generate-csv swagger help
db-up:
	docker compose up -d

help:          
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | awk -F':.*?## ' '{printf "  %-22s %s\n", $$1, $$2}'

init:         
	@test -f Makefile.mk || cp Makefile.mk.dist Makefile.mk
	@test -f .env        || cp .env.dist .env
	@echo "✓ init: Makefile.mk and .env are ready"

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

swagger:
	swag init -g cmd/api/main.go