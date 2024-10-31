package main

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/handlers"
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"github.com/ThamirisMonteiro/rate-limiter/internal/middleware"
	"github.com/ThamirisMonteiro/rate-limiter/internal/ratelimiter"
	"log"
	"net/http"
	"os"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisRepo := database.NewRedisRepository(redisAddr)

	rateLimiterService := ratelimiter.NewService(redisRepo)

	mux := http.NewServeMux()
	mux.Handle("/api/endpoint", middleware.RateLimiterMiddleware(rateLimiterService)(http.HandlerFunc(handlers.Handler)))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
