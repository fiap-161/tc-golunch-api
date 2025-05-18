package postgre_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivens/postgre"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(value any) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query any, args ...any) *gorm.DB {
	mockArgs := m.Called(query, args)
	return mockArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest any, conds ...any) *gorm.DB {
	mockArgs := m.Called(dest, conds)
	return mockArgs.Get(0).(*gorm.DB)
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockDB := new(MockDB)
		repo := postgre.NewRepository(mockDB)

		customer := model.Customer{
			Name:  "John Doe",
			Email: "test@example.com",
			CPF:   "123456",
		}

		mockDB.On("Create", &customer).Return(&gorm.DB{Error: nil})

		saved, err := repo.Create(context.Background(), customer)

		assert.NoError(t, err)
		assert.NotNil(t, saved)
		mockDB.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockDB := new(MockDB)
		repo := postgre.NewRepository(mockDB)

		customer := model.Customer{
			Name:  "John Doe",
			Email: "test@example.com",
			CPF:   "123456",
		}

		expectedErr := errors.New("database error")

		mockDB.On("Create", &customer).Return(&gorm.DB{Error: expectedErr})

		saved, err := repo.Create(context.Background(), customer)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, saved.ID, "")
		mockDB.AssertExpectations(t)
	})
}
