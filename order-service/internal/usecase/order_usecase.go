package usecase

import (
	"context"
	"order-service/internal/domain"
	"order-service/internal/repository"
)

type OrderUsecase struct {
	repo repository.OrderRepository
}

func NewOrderUsecase(repo repository.OrderRepository) *OrderUsecase {
	return &OrderUsecase{repo: repo}
}

func (uc *OrderUsecase) CreateOrder(ctx context.Context, userID string, items []domain.OrderItem) (*domain.Order, error) {
	order := &domain.Order{
		UserID: userID,
		Items:  items,
		Status: "created",
	}

	// Calculate total
	for _, item := range items {
		order.Total += item.Price * float64(item.Quantity)
	}

	if err := uc.repo.Create(ctx, order); err != nil {
		return nil, err
	}
	return order, nil
}

func (uc *OrderUsecase) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	return uc.repo.GetByID(ctx, orderID)
}

func (uc *OrderUsecase) UpdateOrderStatus(ctx context.Context, orderID, status string) error {
	return uc.repo.UpdateStatus(ctx, orderID, status)
}

func (uc *OrderUsecase) ListUserOrders(ctx context.Context, userID string, limit, offset int32) ([]*domain.Order, error) {
	return uc.repo.ListByUser(ctx, userID, limit, offset)
}