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
# DATABASE (một DSN duy nhất, ưu tiên env POSTGRES_DSN)
# ======================
# Ví dụ VPS (dev.yml): postgres://user_postgres:password_postgres@103.90.226.96:5439/gogogo?sslmode=disable
# Trên server PostgreSQL cần cài PostGIS (extension cho bảng zones), ví dụ: apt install postgresql-14-postgis-3

POSTGRES_DSN ?= postgres://user_postgres:password_postgres@103.90.226.96:5439/gogogo?sslmode=disable
MYSQL_DSN    ?= gogin_user:123456@tcp(localhost:3306)/gogin?parseTime=true

# ======================
# GOOSE
# ======================

PG_MIGRATION_DIR = internal/orm/postgres/migrations
PG_SEED_DIR      = internal/orm/postgres/migrations/seeds
PG_SEED_TABLE   = goose_seed_version
MYSQL_MIGRATION_DIR = internal/orm/mysql/migrations

# --- Postgres migrations ---

pg-new:
	@goose -dir $(PG_MIGRATION_DIR) create $(name) sql

pg-up:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" up -allow-missing

pg-down:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" down

pg-status:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" status

pg-reset:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" reset

# --- Postgres seeds ---

pg-seed-up:
	@goose -dir $(PG_SEED_DIR) -table $(PG_SEED_TABLE) postgres "$(POSTGRES_DSN)" up -allow-missing

pg-seed-down:
	@goose -dir $(PG_SEED_DIR) -table $(PG_SEED_TABLE) postgres "$(POSTGRES_DSN)" down

pg-seed-status:
	@goose -dir $(PG_SEED_DIR) -table $(PG_SEED_TABLE) postgres "$(POSTGRES_DSN)" status

pg-seed-new:
	@goose -dir $(PG_SEED_DIR) create $(name) sql

pg-seed: pg-seed-up

# --- MySQL migrations ---

mysql-new:
	@goose -dir $(MYSQL_MIGRATION_DIR) create $(name) sql

mysql-up:
	@goose -dir $(MYSQL_MIGRATION_DIR) mysql "$(MYSQL_DSN)" up

mysql-down:
	@goose -dir $(MYSQL_MIGRATION_DIR) mysql "$(MYSQL_DSN)" down

mysql-status:
	@goose -dir $(MYSQL_MIGRATION_DIR) mysql "$(MYSQL_DSN)" status

mysql-reset:
	@goose -dir $(MYSQL_MIGRATION_DIR) mysql "$(MYSQL_DSN)" reset
