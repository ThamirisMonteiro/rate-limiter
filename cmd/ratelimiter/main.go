package main

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"net/http"
	"os"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")

	database.NewRedisRepository(redisAddr)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
