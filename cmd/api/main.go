package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/Evance365/kozi/internal/config"
    "github.com/Evance365/kozi/internal/db"
    "github.com/Evance365/kozi/internal/handlers"
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

    router := gin.Default()

    router.POST("/results", handlers.PostResults(conn))
    router.GET("/matches", handlers.GetMatches(conn))

    log.Println("kozi running on :8080")
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
