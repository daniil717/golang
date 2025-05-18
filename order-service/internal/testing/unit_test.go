package testing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"order-service/internal/model"
	"order-service/internal/usecase"
)

// 🔧 Mock репо
type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, order *model.Order) (string, error) {
	args := m.Called(ctx, order)
	return args.String(0), args.Error(1)
}

func (m *MockOrderRepo) FindByID(ctx context.Context, id string) (*model.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderRepo) UpdateStatus(ctx context.Context, id, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockOrderRepo) FindByUserID(ctx context.Context, userID string) ([]*model.Order, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Order), args.Error(1)
}

// 🔧 Mock паблишер
type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) PublishOrderCreated(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func getSampleOrder() *model.Order {
	return &model.Order{
		UserID: "user123",
		Products: []model.Product{
			{ProductID: "507f1f77bcf86cd799439011", Quantity: 2},
		},
	}
}

func TestCreateOrder_Success(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	mockPub := new(MockPublisher)
	uc := usecase.NewOrderUsecase(mockRepo, mockPub)

	order := getSampleOrder()
	mockRepo.On("Create", mock.Anything, order).Return("order123", nil)
	mockPub.On("PublishOrderCreated", mock.Anything, order).Return(nil)

	id, err := uc.CreateOrder(context.Background(), order)

	assert.NoError(t, err)
	assert.Equal(t, "order123", id)
	mockRepo.AssertExpectations(t)
	mockPub.AssertExpectations(t)
}

func TestCreateOrder_InvalidInput(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	mockPub := new(MockPublisher)
	uc := usecase.NewOrderUsecase(mockRepo, mockPub)

	order := &model.Order{} // invalid: no UserID or Products

	id, err := uc.CreateOrder(context.Background(), order)

	assert.Error(t, err)
	assert.Empty(t, id)
}

func TestGetOrder(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	mockPub := new(MockPublisher)
	uc := usecase.NewOrderUsecase(mockRepo, mockPub)

	expectedOrder := getSampleOrder()
	expectedOrder.ID = "order123"

	mockRepo.On("FindByID", mock.Anything, "order123").Return(expectedOrder, nil)

	order, err := uc.GetOrder(context.Background(), "order123")

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestUpdateOrderStatus(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	mockPub := new(MockPublisher)
	uc := usecase.NewOrderUsecase(mockRepo, mockPub)

	mockRepo.On("UpdateStatus", mock.Anything, "order123", "COMPLETED").Return(nil)

	err := uc.UpdateOrderStatus(context.Background(), "order123", "COMPLETED")

	assert.NoError(t, err)
}

func TestUpdateOrderStatus_InvalidStatus(t *testing.T) {
	mockRepo := new(MockOrderRepo)
	mockPub := new(MockPublisher)
	uc := usecase.NewOrderUsecase(mockRepo, mockPub)

	err := uc.UpdateOrderStatus(context.Background(), "order123", "SHIPPED") // invalid

	assert.Error(t, err)
}