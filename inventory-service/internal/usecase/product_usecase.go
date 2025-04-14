package usecase

import (
	"context"
	"inventory-servicee/internal/domain"
	"inventory-servicee/internal/repository"
)

type ProductUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo: repo}
}

func (uc *ProductUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	return uc.repo.Create(ctx, product)
}

func (uc *ProductUsecase) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *ProductUsecase) UpdateProduct(ctx context.Context, id string, product *domain.Product) error {
	return uc.repo.Update(ctx, id, product)
}

func (uc *ProductUsecase) DeleteProduct(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *ProductUsecase) ListProducts(ctx context.Context, limit, offset int32) ([]*domain.Product, error) {
	return uc.repo.List(ctx, limit, offset)
}
