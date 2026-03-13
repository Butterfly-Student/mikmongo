package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	_ "mikmongo/internal/migration"
	"mikmongo/internal/seeder"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://" + env("DB_USER", "mikmongo") + ":" +
			env("DB_PASSWORD", "mikmongo") + "@" +
			env("DB_HOST", "localhost") + ":" +
			env("DB_PORT", "5432") + "/" +
			env("DB_NAME", "mikmongo") + "?sslmode=" +
			env("DB_SSLMODE", "disable")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations first
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}
	if err := goose.Up(db, "."); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed")

	// Run seeder
	apiPort, _ := strconv.Atoi(env("SEED_ROUTER_API_PORT", "0"))
	cfg := seeder.Config{
		AdminEmail:    os.Getenv("SEED_ADMIN_EMAIL"),
		AdminPassword: os.Getenv("SEED_ADMIN_PASSWORD"),
		AdminName:     os.Getenv("SEED_ADMIN_NAME"),
		AdminPhone:    os.Getenv("SEED_ADMIN_PHONE"),

		RouterName:     os.Getenv("SEED_ROUTER_NAME"),
		RouterAddress:  os.Getenv("SEED_ROUTER_ADDRESS"),
		RouterAPIPort:  apiPort,
		RouterUsername: os.Getenv("SEED_ROUTER_USERNAME"),
		RouterPassword: os.Getenv("SEED_ROUTER_PASSWORD"),

		EncryptionKey: os.Getenv("JWT_SECRET"),
	}

	if err := seeder.New(db, cfg).Run(context.Background()); err != nil {
		log.Fatalf("Seed failed: %v", err)
	}
	log.Println("Seed completed successfully")
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
