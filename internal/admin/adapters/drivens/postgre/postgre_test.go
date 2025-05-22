package postgre_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivens/postgre"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/model"
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

		admin := model.Admin{
			Email: "test@example.com",
		}

		mockDB.On("Create", &admin).Return(&gorm.DB{Error: nil})

		err := repo.Create(context.Background(), admin)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockDB := new(MockDB)
		repo := postgre.NewRepository(mockDB)

		admin := model.Admin{
			Email: "test@example.com",
		}

		expectedErr := errors.New("database error")

		mockDB.On("Create", &admin).Return(&gorm.DB{Error: expectedErr})

		err := repo.Create(context.Background(), admin)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockDB.AssertExpectations(t)
	})
}
