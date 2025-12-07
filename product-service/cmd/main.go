package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"product-service/internal/config"
	"product-service/internal/handler"
	"product-service/internal/infra"
	"product-service/internal/repository"
	"product-service/internal/router"
	"product-service/internal/service"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()

	infoLog := log.New(log.Writer(), "INFO: ", log.Ldate|log.Ltime)
	errorLog := log.New(log.Writer(), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// DB init via infra package
	db := infra.NewPostgres(cfg, infoLog, errorLog)
	defer db.Close()

	// Wiring DI
	prodRepo := repository.NewPostgresProductRepository(db)
	prodSvc := service.NewProductService(prodRepo)
	prodHandler := handler.NewProductHandler(prodSvc, infoLog, errorLog)
	healthHandler := handler.NewHealthHandler(db)

	r := router.NewRouter(prodHandler, healthHandler)

	// HTTP Server configuration
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		infoLog.Printf("server listening at http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errorLog.Fatalf("server failed to start: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Block until signal received
	<-quit
	infoLog.Println("shutting down server gracefully...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		errorLog.Printf("server forced to shutdown: %v", err)
	}

	infoLog.Println("server stopped")
}
