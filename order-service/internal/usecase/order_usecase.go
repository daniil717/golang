package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	queue "order-service/internal/events"
	"order-service/internal/model"
	"order-service/internal/redis"
	"order-service/internal/repository"
)

type OrderUsecase struct {
	repo      repository.OrderRepository
	publisher queue.Publisher
}

func NewOrderUsecase(repo repository.OrderRepository, publisher queue.Publisher) *OrderUsecase {
	return &OrderUsecase{
		repo:      repo,
		publisher: publisher,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	if order.UserID == "" || len(order.Products) == 0 {
		return "", errors.New("invalid order data")
	}
	if order.Status == "" {
		order.Status = "PENDING"
	}

	for _, product := range order.Products {
		if product.ProductID == "" {
			return "", errors.New("invalid product ID in order")
		}
		if len(product.ProductID) != 24 {
			log.Printf("Warning: ProductID %s does not match expected ObjectID format", product.ProductID)
		}
	}

	id, err := u.repo.Create(ctx, order)
	if err != nil {
		return "", err
	}
	order.ID = id

	// ✅ NATS publisher бар ма, соны тексер
	if u.publisher != nil {
		err = u.publisher.PublishOrderCreated(ctx, order)
		if err != nil {
			log.Printf("[NATS] failed to publish order.created event: %v", err)
		}
	}

	return id, nil
}

func (u *OrderUsecase) GetOrder(ctx context.Context, id string) (*model.Order, error) {
	if id == "" {
		return nil, errors.New("invalid order ID")
	}
	return u.repo.FindByID(ctx, id)
}

func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, id string, status string) error {
	if id == "" || status == "" {
		return errors.New("invalid input data")
	}
	validStatuses := map[string]bool{
		"PENDING":   true,
		"COMPLETED": true,
		"CANCELLED": true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}
	return u.repo.UpdateStatus(ctx, id, status)
}

func (u *OrderUsecase) ListUserOrders(ctx context.Context, userID string) ([]*model.Order, error) {
	if userID == "" {
		return nil, errors.New("invalid user ID")
	}

	// Cache key for orders
	cacheKey := fmt.Sprintf("orders:user:%s", userID)

	// Check the cache first
	cachedOrders, err := redis.GetFromCache[[]*model.Order](cacheKey)
	if err != nil {
		return nil, err
	}

	if cachedOrders != nil {
		// If we have cached data, return it
		return *cachedOrders, nil
	}

	// Otherwise, fetch from the repository
	orders, err := u.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Save the result in cache for subsequent requests
	_ = redis.SetToCache(cacheKey, orders, 10*time.Minute)

	return orders, nil
}
