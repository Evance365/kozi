package main

import (
    "log"
    "os"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/Evance365/kozi/internal/config"
    "github.com/Evance365/kozi/internal/db"
    "github.com/Evance365/kozi/internal/handlers"
)

func main() {
    godotenv.Load()

    cfg := config.Load()

    conn, err := db.Connect(cfg.DSN())
    if err != nil {
        log.Fatalf("database error: %v", err)
    }
    defer conn.Close()

    router := gin.Default()

    allowedOrigins := []string{"http://localhost:3000"}
    if origin := os.Getenv("ALLOWED_ORIGIN"); origin != "" {
        allowedOrigins = append(allowedOrigins, origin)
    }

    router.Use(cors.New(cors.Config{
        AllowOrigins:     allowedOrigins,
        AllowMethods:     []string{"GET", "POST", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type"},
        AllowCredentials: false,
    }))

    router.POST("/results", handlers.PostResults(conn))
    router.GET("/matches", handlers.GetMatches(conn))

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("kozi running on :%s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("server error: %v", err)
    }
}
