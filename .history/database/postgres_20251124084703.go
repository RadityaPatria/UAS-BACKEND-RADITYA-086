package database

import (
	"context"
	"fmt"
	"log"

	"UAS-backend/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectPostgres(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Failed to connect PostgreSQL:", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("PostgreSQL ping failed:", err)
	}

	DB = pool
	log.Println("Connected to PostgreSQL successfully using pgx.")
}
