package main

import (
	"log"
	"net/http"
	"product-service/internal/config"
	"product-service/internal/handler"
	"product-service/internal/infra"
	"product-service/internal/repository"
	"product-service/internal/router"
	"product-service/internal/service"
)

func main() {
	cfg := config.Load()

	infoLog := log.New(log.Writer(), "INFO: ", log.Ldate|log.Ltime)
	errorLog := log.New(log.Writer(), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// DB init via infra package
	db := infra.NewPostgres(cfg, infoLog, errorLog)

	// wiring DI
	prodRepo := repository.NewPostgresProductRepository(db)
	prodSvc := service.NewProductService(prodRepo)
	prodHandler := handler.NewProductHandler(prodSvc, infoLog, errorLog)

	r := router.NewRouter(prodHandler)

	infoLog.Printf("server listening at http://localhost:%s/products", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		errorLog.Fatalf("server failed: %v", err)
	}
}
