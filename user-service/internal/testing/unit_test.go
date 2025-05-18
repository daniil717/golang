package testing

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "user-service/internal/model"
    "user-service/internal/usecase"
)

// 🔧 Mock репозиторий
type MockUserRepository struct {
    mock.Mock
}

// Cleanup әдісі, егер интерфейсте бар болса (егер жоқ болса, алып таста)
func (m *MockUserRepository) Cleanup(ctx context.Context) error {
    return nil
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) (string, error) {
    args := m.Called(ctx, user)
    return args.String(0), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
    args := m.Called(ctx, username)
    user := args.Get(0)
    if user == nil {
        return nil, args.Error(1)
    }
    return user.(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
    args := m.Called(ctx, id)
    user := args.Get(0)
    if user == nil {
        return nil, args.Error(1)
    }
    return user.(*model.User), args.Error(1)
}


// 🧪 Unit test
func TestCreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    uc := usecase.NewUserUsecase(mockRepo)

    testUser := &model.User{
        ID:       "1",
        Username: "Alice",
        Email:    "alice@example.com",
        Password: "password123",
    }

    // Моктың Create әдісі шақырылғанда, параметрлерге сай жұмыс істеуін орнатамыз
    mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(user *model.User) bool {
        return user.Username == testUser.Username && user.Email == testUser.Email
    })).Return("1", nil)

    id, err := uc.CreateUser(context.Background(), testUser)

    assert.NoError(t, err)
    assert.Equal(t, "1", id)
    mockRepo.AssertExpectations(t)
}
