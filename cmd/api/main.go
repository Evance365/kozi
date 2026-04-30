package main

import (
    "log"

    "github.com/joho/godotenv"
    "github.com/Evance365/kozi/internal/config"
    "github.com/Evance365/kozi/internal/db"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("failed to load .env file")
    }

    cfg := config.Load()

    conn, err := db.Connect(cfg.DSN())
    if err != nil {
        log.Fatalf("database error: %v", err)
    }
    defer conn.Close()

    log.Println("kozi is ready")
}
