package main

import (
	"database/sql"
	"log"

	"weight-loss-challenge/internal/database"
	"weight-loss-challenge/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
	// _ "github.com/mattn/go-sqlite3" //local sqlite
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	// // local sqlite
	// db, err := sql.Open("sqlite3", "./data.db")

	dsn := env.GetEnvString("TURSO_DATABASE_URL", "") + "?authToken=" + env.GetEnvString("TURSO_AUTH_TOKEN", "")

	db, err := sql.Open("libsql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
