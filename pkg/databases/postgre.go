package databases

import (
	"context"
	"log"
	"sport-assistance/pkg/configs"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(cfg *configs.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseConfig.DBConnectionString)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v", err)
	}

	poolConfig.MaxConns = int32(cfg.DatabaseConfig.DBMaxConn)
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Printf("Connected to PostgreSQL database (max pool size: %d)\n", poolConfig.MaxConns)

	return pool, nil
}
