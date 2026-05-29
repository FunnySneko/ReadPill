include .env
export

.PHONY: migrate-up
migrate-up:
	migrate -path ./migrations -database postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE} up
.PHONY: migrate-down
migrate-down:
	migrate -path ./migrations -database postgresql://${PGUSER}:${PGPASSWORD}@${PGHOST}:${PGPORT}/${PGDATABASE} down