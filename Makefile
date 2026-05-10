run:
	POSTGRES_USER=jigoku POSTGRES_PASSWORD=123 POSTGRES_DB=students_db POSTGRES_HOST=localhost POSTGRES_PORT=5433 go run ./cmd/api

db-up:
	docker compose up -d

db-down:
	docker compose down
	
generate-csv:
	python3 scripts/generate_students_csv.py

import-csv:
	curl -X POST http://localhost:8080/students/import -F "file=@raw_data/students.csv"

run-distribution:
	curl -X POST http://localhost:8080/distribution/run

export-distribution:
	curl http://localhost:8080/distribution/export -o raw_data/distribution_export.csv