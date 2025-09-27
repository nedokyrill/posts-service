package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nedokyrill/posts-service/pkg/logger"
)

func Connect() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		logger.Logger.Errorf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	logger.Logger.Info("Connected to Postgres successfully")
	return conn, nil
}
