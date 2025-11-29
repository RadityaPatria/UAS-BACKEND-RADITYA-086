package database

import (
	"context"
	"fmt"
	"log"

	"UAS-backend/config"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectPostgres(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect PostgreSQL:", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatal("❌ PostgreSQL ping failed:", err)
	}

	DB = conn
	log.Println("✅ Connected to PostgreSQL successfully (pgx.Conn).")
}
