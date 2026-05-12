package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NickSFU/shortlink-service/internal/cache"
	"github.com/NickSFU/shortlink-service/internal/config"
	"github.com/NickSFU/shortlink-service/internal/db"
	"github.com/NickSFU/shortlink-service/internal/router"
)

func Run() error {
	cfg := config.Load()
	pool, err := db.NewPostgres(cfg)

	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer pool.Close()
	redisClient := cache.NewRedis(cfg)

	if err := cache.Ping(redisClient); err != nil {
		log.Fatalf("redis error: %v", err)
	}
	router := router.NewRouter(pool, redisClient)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		fmt.Printf("server started on port %s\n", cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
