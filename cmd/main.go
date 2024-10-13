package main

import (
	"log"

	"github.com/MXkodo/cash-server/config"
	"github.com/MXkodo/cash-server/internal/app"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := app.Run(cfg, rdb); err != nil {
		log.Fatal("Failed to start app: ", err)
	}
}
