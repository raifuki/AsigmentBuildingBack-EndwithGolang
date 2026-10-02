package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/raifuki/task-management/internal/config"
	"github.com/raifuki/task-management/internal/database"
	"github.com/raifuki/task-management/internal/routes"
)

func main() {
	cfg := config.Load()
	db := database.NewPostgres(cfg)
	rdb := database.NewRedis(cfg)

	// Cache layer
	// cacheLayer := cache.NewRedisCache(rdb) // sẽ dùng trong routes

	router := routes.SetupRouter(cfg, db, rdb)

	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.AppPort
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("🚀 Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("✅ Server exited")
}
