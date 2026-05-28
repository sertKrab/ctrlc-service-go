package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[Migrate] config: %v", err)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Fatalf("[Migrate] init: %v", err)
	}
	defer m.Close()

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("[Migrate] up: %v", err)
		}
		log.Println("[Migrate] up: done")
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("[Migrate] down: %v", err)
		}
		log.Println("[Migrate] down: 1 step rolled back")
	case "drop":
		if err := m.Drop(); err != nil {
			log.Fatalf("[Migrate] drop: %v", err)
		}
		log.Println("[Migrate] drop: done")
	default:
		log.Fatalf("[Migrate] unknown command: %s (use: up | down | drop)", command)
	}
}
