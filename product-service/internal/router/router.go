package router

import (
	"product-service/internal/handler"
	"product-service/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(prodHandler *handler.ProductHandler) *chi.Mux {
	r := chi.NewRouter()

	// middleware global
	r.Use(middleware.Logging)
	r.Use(middleware.CORS)
	r.Use(middleware.RateLimit)

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
