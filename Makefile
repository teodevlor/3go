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

POSTGRES_DSN=postgres://gogin_user:123456@localhost:5432/gogin?sslmode=disable
MYSQL_DSN=gogin_user:123456@tcp(localhost:3306)/gogin?parseTime=true

# ======================
# GOOSE DIRECTORIES
# ======================

PG_MIGRATION_DIR=internal/orm/postgres/migrations
MYSQL_MIGRATION_DIR=internal/orm/mysql/migrations

# ======================
# POSTGRES MIGRATIONS
# ======================

# make pg-new name=add_users_table
pg-new:
	@goose -dir $(PG_MIGRATION_DIR) create $(name) sql

pg-up:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" up

pg-down:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" down

pg-reset:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" reset

pg-status:
	@goose -dir $(PG_MIGRATION_DIR) postgres "$(POSTGRES_DSN)" status

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

