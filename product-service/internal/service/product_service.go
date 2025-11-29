package service

import (
	"context"
	"product-service/internal/models"
	"product-service/internal/repository"
)

type ProductService interface {
	Create(ctx context.Context, p *models.Product) error
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	List(ctx context.Context) ([]*models.Product, error)
	ListPaged(ctx context.Context, limit, offset int) ([]*models.Product, error)
	Search(ctx context.Context, q string) ([]*models.Product, error)
	Update(ctx context.Context, p *models.Product) (*models.Product, error)
	Patch(ctx context.Context, id int64, patch *models.ProductPatch) (*models.Product, error)
	Delete(ctx context.Context, id int64) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

// // @SERVICE:CREATE-BEGIN
func (s *productService) Create(ctx context.Context, p *models.Product) error {
	return s.repo.Create(ctx, p)
}

//// @SERVICE:CREATE-END

// // @SERVICE:READ-BEGIN
func (s *productService) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) List(ctx context.Context) ([]*models.Product, error) {
	return s.repo.List(ctx)
}

func (s *productService) ListPaged(ctx context.Context, page, limit int) ([]*models.Product, error) {
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	return s.repo.ListPaged(ctx, limit, offset)
}

func (s *productService) Search(ctx context.Context, q string) ([]*models.Product, error) {
	return s.repo.Search(ctx, q)
}

//// @SERVICE:READ-END

// // @SERVICE:UPDATE-BEGIN
func (s *productService) Update(ctx context.Context, p *models.Product) (*models.Product, error) {
	return s.repo.Update(ctx, p)
}

func (s *productService) Patch(ctx context.Context, id int64, patch *models.ProductPatch) (*models.Product, error) {
	return s.repo.Patch(ctx, id, patch)
}

//// @SERVICE:UPDATE-END

// // @SERVICE:DELETE-BEGIN
func (s *productService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

//// @SERVICE:DELETE-END
