package infra

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"product-service/internal/config"

	_ "github.com/lib/pq"
)

func NewPostgres(cfg *config.Config, infoLog, errorLog *log.Logger) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		errorLog.Fatalf("failed to open db: %v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		errorLog.Fatalf("db ping failed: %v", err)
	}

	// Connection Pool Settings (safe default)
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	infoLog.Println("connected to postgres")

	return db
}
