package main

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/handlers"
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"github.com/ThamirisMonteiro/rate-limiter/internal/middleware"
	"github.com/ThamirisMonteiro/rate-limiter/internal/ratelimiter"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	redisAddr, reqLimit, blockTimeIP, blockTimeToken := getEnvVariables()

	redisRepo := database.NewRedisRepository(redisAddr)

	rateLimiterService, err := ratelimiter.NewService(redisRepo, reqLimit, time.Duration(blockTimeIP)*time.Second, time.Duration(blockTimeToken)*time.Second)
	if err != nil {
		log.Fatalf("Erro ao criar o serviço de rate limiter: %v", err)
	}

	http.Handle("/", middleware.RateLimiterMiddleware(rateLimiterService)(http.HandlerFunc(handlers.Handler)))

	log.Println("Servidor iniciado na porta 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

func getEnvVariables() (string, int, int, int) {
	redisAddr := os.Getenv("REDIS_ADDR")
	reqLimitStr := os.Getenv("REQ_LIMIT")
	blockTimeIPStr := os.Getenv("BLOCK_TIME_IP")
	blockTimeTokenStr := os.Getenv("BLOCK_TIME_TOKEN")

	reqLimit, _ := strconv.Atoi(reqLimitStr)
	blockTimeIP, _ := strconv.Atoi(blockTimeIPStr)
	blockTimeToken, _ := strconv.Atoi(blockTimeTokenStr)

	if reqLimit <= 0 || blockTimeIP <= 0 || blockTimeToken <= 0 {
		log.Fatalf("Valores inválidos para limites: REQ_LIMIT: %d, BLOCK_TIME_IP: %d, BLOCK_TIME_TOKEN: %d", reqLimit, blockTimeIP, blockTimeToken)
	}

	return redisAddr, reqLimit, blockTimeIP, blockTimeToken
}
