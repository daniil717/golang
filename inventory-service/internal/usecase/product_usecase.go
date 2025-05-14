package usecase

import (
	"context"
	"errors"
	"fmt"
	"inventory-service/internal/model"
	"inventory-service/internal/redis"
	"inventory-service/internal/repository"
	"time"
)

type ProductUsecase struct {
    repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) *ProductUsecase {
    return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, p *model.Product) (string, error) {
    if p.Name == "" {
        return "", errors.New("name is required")
    }
    if p.Description == "" {
        return "", errors.New("description is required")
    }
    if p.Category == "" {
        return "", errors.New("category is required")
    }
    if p.Stock < 0 {
        return "", errors.New("stock cannot be negative")
    }
    if p.Price < 0 {
        return "", errors.New("price cannot be negative")
    }
    return u.repo.Create(ctx, p)
}

func (u *ProductUsecase) GetProduct(ctx context.Context, id string) (*model.Product, error) {
    if id == "" {
        return nil, errors.New("id is required")
    }

    key := fmt.Sprintf("product:%s", id)

    // 1. Redis кэштен іздеу
    cached, err := redis.GetFromCache[model.Product](key)
    if err != nil {
        return nil, err
    }
    if cached != nil {
        return cached, nil
    }

    // 2. Базадан алу
    product, err := u.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. Redis-ке сақтау
    _ = redis.SetToCache(key, product, time.Hour)

    return product, nil
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, p *model.Product) error {
    if p.ID == "" {
        return errors.New("id is required")
    }
    if p.Name == "" || p.Description == "" || p.Category == "" {
        return errors.New("name, description and category are required")
    }
    if p.Stock < 0 {
        return errors.New("stock cannot be negative")
    }
    if p.Price < 0 {
        return errors.New("price cannot be negative")
    }

    err := u.repo.Update(ctx, p)
    if err != nil {
        return err
    }

    // Кэшті өшіру
    key := fmt.Sprintf("product:%s", p.ID)
    _ = redis.DeleteCache(key)

    return nil
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id string) error {
    if id == "" {
        return errors.New("id is required")
    }
    return u.repo.Delete(ctx, id)
}

func (u *ProductUsecase) ListProducts(ctx context.Context, category string, page, limit int32) ([]*model.Product, error) {
    // Әдепкі мәндер
    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10 // Мұны өзіңе ыңғайлы мәнге өзгертуге болады
    }

    return u.repo.List(ctx, category, page, limit)
}

func (u *ProductUsecase) DecreaseStock(ctx context.Context, productID string, quantity int32) error {
    return u.repo.DecreaseStock(ctx, productID, quantity)
}

