.PHONY: create-migrate, process-migration

create-migrate:
	migrate create -ext sql -seq -dir migrations 'init_schema'

process-migration:
	migrate -path migrations -database 'postgres://postgres:password@localhost:5432/storage?sslmode=disable' up
