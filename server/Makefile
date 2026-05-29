include .env
export

.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations -database postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=disable up

.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations -database postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE}?sslmode=disable down

.PHONY: setup
setup:
	@set -a && . ./.env && set +a && \
	PGPASSWORD=$$PGPASSWORD psql -h $$PGHOST -p $$PGPORT -U $$PGUSER -tc "SELECT 1 FROM pg_database WHERE datname='$$PGDATABASE'" | grep -q 1 || \
	PGPASSWORD=$$PGPASSWORD psql -h $$PGHOST -p $$PGPORT -U $$PGUSER -c "CREATE DATABASE $$PGDATABASE" && \
	migrate -path ./migrations -database "postgres://$$PGUSER:$$PGPASSWORD@$$PGHOST:$$PGPORT/$$PGDATABASE?sslmode=disable" up