package db

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

func Connect(dsn string) (*sql.DB, error) {
    conn, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open db connection: %w", err)
    }

    if err := conn.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    conn.SetMaxOpenConns(25)
    conn.SetMaxIdleConns(5)

    log.Println("database connection established")
    return conn, nil
}
