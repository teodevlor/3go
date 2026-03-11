package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN is required")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	isSeed := len(os.Args) > 1 && os.Args[1] == "seed"

	if isSeed {
		goose.SetTableName("goose_db_version_seeds")
		if err := goose.Up(db, "orm/postgres/migrations/seeds"); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Seed completed successfully")
	} else {
		if err := goose.Up(db, "orm/postgres/migrations"); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Migration completed successfully")
	}
}
