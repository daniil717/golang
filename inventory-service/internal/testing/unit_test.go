package testing

import (
	"context"
	"testing"

	"inventory-service/internal/model"
	"inventory-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock репозиторий интерфейсі
type MockProductRepo struct {
	mock.Mock
}

func (m *MockProductRepo) Create(ctx context.Context, p *model.Product) (string, error) {
	args := m.Called(ctx, p)
	return args.String(0), args.Error(1)
}

func (m *MockProductRepo) GetByID(ctx context.Context, id string) (*model.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepo) Update(ctx context.Context, p *model.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProductRepo) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepo) List(ctx context.Context, category string, page, limit int32) ([]*model.Product, error) {
	args := m.Called(ctx, category, page, limit)
	return args.Get(0).([]*model.Product), args.Error(1)
}

func (m *MockProductRepo) DecreaseStock(ctx context.Context, productID string, quantity int32) error {
	args := m.Called(ctx, productID, quantity)
	return args.Error(0)
}

// Тесттер

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepo)
	uc := usecase.NewProductUsecase(mockRepo)

	p := &model.Product{
		Name:        "Product1",
		Description: "Desc",
		Category:    "Cat",
		Stock:       10,
		Price:       100,
	}

	mockRepo.On("Create", mock.Anything, p).Return("productID123", nil)

	id, err := uc.CreateProduct(context.Background(), p)

	assert.NoError(t, err)
	assert.Equal(t, "productID123", id)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_InvalidInput(t *testing.T) {
	uc := usecase.NewProductUsecase(nil)

	_, err := uc.CreateProduct(context.Background(), &model.Product{Name: ""})
	assert.Error(t, err)

	_, err = uc.CreateProduct(context.Background(), &model.Product{Name: "Valid", Description: "", Category: "Cat"})
	assert.Error(t, err)

	_, err = uc.CreateProduct(context.Background(), &model.Product{Name: "Valid", Description: "Desc", Category: "Cat", Stock: -1})
	assert.Error(t, err)
}
