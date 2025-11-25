DB_URL ?= postgresql://root:secret@192.168.29.20:5432/authentication?sslmode=disable

# setting up docker, database and migration
# ******************************************************************************** 
create-auth-container:
	fuser -k 5432/tcp 2>/dev/null || true && docker run --name auth-package -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine

drop-auth-container:
	docker stop auth-package
	docker rm auth-package

psql-shell:
	docker exec -it auth-package psql

create-db:
	docker exec -it auth-package createdb --username=root --owner=root authentication

migrate-up:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path internal/db/migration -database "$(DB_URL)" -verbose down

migrate-fresh: migratedown migrateup
	@echo "Fresh migration complete"
# ********************************************************************************

# setup for sqlc
# ********************************************************************************
sqlc:
	sqlc generate
# ********************************************************************************

db-test:
	go test -v -cover -count=1 ./internal/db/tests

serivce-test:
	go test -v -cover -count=1 ./internal/services/tests

# running server  HACK: Remove it after production
server:
	fuser -k 8000/tcp 2>/dev/null || true && go run ./cmd/server/main.go

mock:
	mockgen -source=./internal/db/auth.go -destination=./internal/db/mock/auth.go -package=mock Auth

.PHONY: create-auth-container drop-auth-container psql-shell migrateup migratedown
