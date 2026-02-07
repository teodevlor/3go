run:
	go run cmd/server/main.go

dev:
	air

telegram-test:
	go run cmd/init/telegram.go

SQLC=sqlc

.PHONY: sqlc-postgres sqlc-mysql sqlc-all sqlc-clean

sqlc-postgres:
	$(SQLC) generate -f sqlc.postgres.yaml

sqlc-mysql:
	$(SQLC) generate -f sqlc.mysql.yaml

sqlc-all: sqlc-postgres sqlc-mysql

sqlc-clean:
	rm -rf internal/orm/db/postgres
	rm -rf internal/orm/db/mysql

# ======================
# DATABASE CONNECTIONS
# ======================

APP_ENV ?= local

POSTGRES_DSN_LOCAL  = postgres://postgres:postgres@localhost:5433/gogogo?sslmode=disable
POSTGRES_DSN_DOCKER = postgres://postgres:postgres@localhost:5433/gogogo?sslmode=disable

MYSQL_DSN_LOCAL  = gogin_user:123456@tcp(localhost:3306)/gogin?parseTime=true
MYSQL_DSN_DOCKER = gogin_user:123456@tcp(mysql:3306)/gogin?parseTime=true

ifeq ($(APP_ENV),docker)
	POSTGRES_DSN = $(POSTGRES_DSN_DOCKER)
	MYSQL_DSN    = $(MYSQL_DSN_DOCKER)
else
	POSTGRES_DSN = $(POSTGRES_DSN_LOCAL)
	MYSQL_DSN    = $(MYSQL_DSN_LOCAL)
endif

# ======================
# GOOSE DIRECTORIES
# ======================

PG_MIGRATION_DIR=internal/orm/postgres/migrations
PG_SEED_DIR=internal/orm/postgres/migrations/seeds
MYSQL_MIGRATION_DIR=internal/orm/mysql/migrations

# ======================
# POSTGRES MIGRATIONS
# ======================

pg-new:
	@goose -dir $(PG_MIGRATION_DIR) create $(name) sql

pg-up:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" up -allow-missing

pg-down:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" down

pg-reset:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" reset

pg-status:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" status

# ======================
# POSTGRES SEEDS
# ======================

pg-seed-up:
	@goose -dir $(PG_SEED_DIR) postgres "$(POSTGRES_DSN)" up

pg-seed-down:
	@goose -dir $(PG_SEED_DIR) postgres "$(POSTGRES_DSN)" down

pg-seed-status:
	@goose -dir $(PG_SEED_DIR) postgres "$(POSTGRES_DSN)" status

pg-seed-new:
	@goose -dir $(PG_SEED_DIR) create $(name) sql

pg-seed: pg-seed-up

# ======================
# MYSQL MIGRATIONS
# ======================

mysql-new:
	@goose -dir $(MYSQL_MIGRATION_DIR) create $(name) sql

mysql-up:
	@goose -dir $(MYSQL_MIGRATION_DIR) mysql "$(MYSQL_DSN)" up

mysql-down:
	@goose -dir $(MYSQL_MIGRATION_DIR) mysql "$(MYSQL_DSN)" down

mysql-reset:
	@goose -dir $(MYSQL_MIGRATION_DIR) mysql "$(MYSQL_DSN)" reset

mysql-status:
	@goose -dir $(MYSQL_MIGRATION_DIR) mysql "$(MYSQL_DSN)" status

# ======================
# DOCKER
# ======================

run-docker:
	chmod +x scripts/start_docker.sh
	./scripts/start_docker.sh

run-goose-up-docker:
	APP_ENV=docker make pg-up 