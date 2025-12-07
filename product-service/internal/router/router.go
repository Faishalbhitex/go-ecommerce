package router

import (
	"product-service/internal/handler"
	"product-service/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(prodHandler *handler.ProductHandler, healthHandler *handler.HealthHandler) *chi.Mux {
	r := chi.NewRouter()

	// middleware global
	r.Use(middleware.Logging)
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimit)

	// Health check endpoints (no rate limiting)
	r.Group(func(r chi.Router) {
		r.Get("/health", healthHandler.Health)
		r.Get("/health/ready", healthHandler.Ready)
		r.Get("/health/live", healthHandler.Live)
	})

	// Product endpoints
	r.Route("/products", func(r chi.Router) {
		r.Post("/", prodHandler.Create)
		r.Get("/", prodHandler.List)

		r.Get("/paged", prodHandler.ListPaged)
		r.Get("/search", prodHandler.Search)

		r.Get("/{id}", prodHandler.GetByID)
		r.Put("/{id}", prodHandler.Update)
		r.Patch("/{id}", prodHandler.Patch)
		r.Delete("/{id}", prodHandler.Delete)
	})

	return r
}
