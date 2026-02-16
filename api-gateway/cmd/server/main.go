package main

import (
    "log"
    "os"

    "github.com/Sacs616/streaming-app/api-gateway/internal"
)

func main() {
    gateway, err := internal.NewGateway()
    if err != nil {
        log.Fatalf("Failed to create gateway: %v", err)
    }

    router := gateway.SetupRouter()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("API Gateway listening on port %s", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
