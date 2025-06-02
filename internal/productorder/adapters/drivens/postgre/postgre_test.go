package postgre_test

import (
	"context"
	"errors"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/adapters/drivens/postgre"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/core/model"
)

type MockDB struct {
	mock.Mock
}

// TODO refactor, move all mocks to a common package so we can reuse if wanted

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
	type args struct {
		context context.Context
		orders  []model.ProductOrder
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Given a valid order, when CreateBulk is called, then it should return the number of created orders",
			args: args{
				orders: []model.ProductOrder{
					{
						Entity: entity.Entity{
							ID: "order_id",
						},
						ProductID: "product_id",
						OrderID:   "order_id",
						Quantity:  1,
						UnitPrice: 100,
					},
					{
						Entity: entity.Entity{
							ID: "order_id_02",
						},
						ProductID: "product_id_02",
						OrderID:   "order_id_02",
						Quantity:  2,
						UnitPrice: 20,
					},
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "Given an error during order creation, when CreateBulk is called, then it should return an error",
			args: args{
				orders: []model.ProductOrder{
					{
						Entity: entity.Entity{
							ID: "order_id",
						},
						ProductID: "product_id",
						OrderID:   "order_id",
						Quantity:  1,
						UnitPrice: 100,
					},
					{
						Entity: entity.Entity{
							ID: "order_id_02",
						},
						ProductID: "product_id_02",
						OrderID:   "order_id_02",
						Quantity:  2,
						UnitPrice: 20,
					},
				},
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			if tt.wantErr {
				mockDB.On("Create", mock.Anything).Return(&gorm.DB{Error: errors.New("error creating orders")})
			} else {
				mockDB.On("Create", mock.Anything).Return(&gorm.DB{RowsAffected: int64(tt.want)})
			}

			repo := postgre.New(mockDB)

			got, err := repo.CreateBulk(tt.args.context, tt.args.orders)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}
