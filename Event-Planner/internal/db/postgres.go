package db

import (
    "context"
    "fmt"
    "os"
    "github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Pool, error) {
    dbURL := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

    pool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        return nil, fmt.Errorf("cannot connect to database: %w", err)
    }

    return pool, nil
}
