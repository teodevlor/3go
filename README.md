# API Backend for Gogogo Delivery System

Golang API using the Gin framework with a layered structure: **API → Controller → Usecase → Transformer → Repository → Mapper → Model → ORM (SQLC)**.

---

## Table of Contents

- [Requirements](#requirements)
- [Project Structure](#project-structure)
- [How to run](#how-to-run)
- [Run with Docker (Makefile flow)](#run-with-docker-makefile-flow)
- [Run locally without Docker](#run-locally-without-docker)
- [Makefile – Command reference](#makefile--command-reference)
- [Database & Migration](#database--migration)

---

## Requirements

- **Go** 1.22+ (or as in `go.mod`)
- **Docker** & **Docker Compose** (for Docker-based run)
- **goose** (migrations): `go install github.com/pressly/goose/v3/cmd/goose@latest`
- **sqlc** (if editing queries/schema): `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- **air** (dev hot-reload, optional): `go install github.com/air-verse/air@latest`

---

## Project Structure

Layered layout: entrypoint `cmd/server/main.go`, config in `config/`, API in `internal/api/v1/`, DI in `internal/registry/`, DB layer in `internal/orm/` (SQLC + goose) and `internal/repository/`. Run config: `config/dev.yml`; Docker: `docker/`, script `scripts/start_docker.sh`.

---

## How to run

**Docker (recommended):**

```bash
make run-docker
make run-goose-up-docker
# API: http://localhost:8080
```

**Local:** Postgres on port 5433, set `config/dev.yml` → `db_postgre.host: localhost`, `db_postgre.port: 5433`. Then:

```bash
make pg-up
make run
# or: make dev  (with air hot-reload)
```

---

## Run with Docker (Makefile flow)

### Step 1: Start Docker

```bash
make run-docker
```

- Makes `scripts/start_docker.sh` executable and runs it.
- Script creates required directories (e.g. `scripts/postgres`) and runs:
  - `docker compose -f docker/docker-compose.yml up -d --build`
- Result:
  - **PostgreSQL** (PostGIS): host port **5433** → 5432 in container; DB `gogogo`, user/pass `postgres`/`postgres`.
  - **API (gogogo-api)**: built from `docker/Dockerfile`, port **8080**, mounts `../config` and `../storage`.

### Step 2: Run migrations (Postgres)

From the **host** (with `goose` installed):

```bash
make run-goose-up-docker
```

- Equivalent to: `APP_ENV=docker make pg-up`.
- `APP_ENV=docker` makes the Makefile use `POSTGRES_DSN` with **localhost:5433**.
- Migrations: `internal/orm/postgres/migrations/`.
- Seeds: `internal/orm/postgres/migrations/seeds/`.

### Step 3: Use the API

- API listens at **http://localhost:8080**.
- Config is read from `config/dev.yml` (mounted in the container).

### Docker flow summary

```
make run-docker
    → scripts/start_docker.sh
        → docker compose -f docker/docker-compose.yml up -d --build
            → postgres (5433:5432)
            → gogogo-api (8080:8080, config/, storage/ mounted)

make run-goose-up-docker
    → APP_ENV=docker make pg-up
        → goose ... postgres "DSN localhost:5433" up
```

---

## Run locally without Docker

1. Install and run Postgres (e.g. local port 5433, DB `gogogo`).
2. Edit `config/dev.yml`: `db_postgre.host` = `localhost`, `db_postgre.port` = `5433` (or your port).
3. Run migrations:

   ```bash
   make pg-up
   ```

   (Default `APP_ENV=local` uses `POSTGRES_DSN_LOCAL` = localhost:5433.)

4. Start the app:

   ```bash
   make run
   ```

   Or with hot-reload (air):

   ```bash
   make dev
   ```

---

## Makefile – Command reference

| Command                      | Description                                              |
|-----------------------------|----------------------------------------------------------|
| `make run`                  | Run app: `go run cmd/server/main.go`                     |
| `make dev`                  | Run with **air** (hot-reload); config in `.air.toml`     |
| `make run-docker`           | Start Docker: script + `docker compose up -d --build`    |
| `make run-goose-up-docker`  | Run Postgres migrations with Docker (`APP_ENV=docker make pg-up`) |
| `make pg-up`                | Run Postgres migrations (DSN depends on APP_ENV)        |
| `make pg-down`              | Rollback one migration                                  |
| `make pg-status`            | Show migration status                                    |
| `make pg-reset`             | Reset migrations (down then up)                          |
| `make pg-seed` / `pg-seed-up` | Run seeds (goose from seeds directory)                |
| `make pg-new name=xxx`      | Create new Postgres migration                            |
| `make sqlc-postgres`        | Generate code from SQLC (Postgres)                        |

Environment:

- **APP_ENV**: `local` (default) or `docker`.
  - `local`: Postgres/MySQL DSN use localhost (e.g. 5433 for Postgres).
  - `docker`: Postgres DSN uses `localhost:5433` for running goose from host into the container.

---

## Database & Migration

- **Engine:** PostgreSQL (PostGIS in Docker).
- **Queries:** SQLC (generated from `internal/orm/postgres/queries/`, config `sqlc.postgres.yaml`).
- **Migrations:** **goose**
  - Migrations: `internal/orm/postgres/migrations/`
  - Seeds: `internal/orm/postgres/migrations/seeds/`
- After changing schema or queries:
  1. Update `.sql` in `queries/` or add a migration.
  2. Run `make sqlc-postgres` if SQLC queries/schema changed.
  3. Run `make pg-up` (local) or `make run-goose-up-docker` (Docker).

Config: **`config/dev.yml`** (with Docker this file is mounted from the host).
