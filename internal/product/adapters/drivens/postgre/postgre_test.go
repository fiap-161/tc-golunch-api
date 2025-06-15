package postgre_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivens/postgre"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Product{})
	require.NoError(t, err)

	return db
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		product  model.Product
		wantErr  string
		validate func(t *testing.T, saved model.Product)
	}{
		{
			name: "success - inserts valid product without error",
			product: model.Product{
				Name:          "Cheeseburger",
				Price:         24.90,
				Description:   "Delicioso cheeseburger com duas carnes.",
				PreparingTime: 15,
				Category:      enum.Meal,
				ImageURL:      "https://example.com/images/cheeseburger.png",
			}.Build(),
			wantErr: "",
			validate: func(t *testing.T, saved model.Product) {
				assert.NotEmpty(t, saved.ID)
				assert.Equal(t, "Cheeseburger", saved.Name)
			},
		},
		{
			name:    "error - missing required fields",
			product: model.Product{},
			wantErr: "name is required",
			validate: func(t *testing.T, saved model.Product) {
				assert.Empty(t, saved.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := postgre.New(db)

			saved, err := repo.Create(context.Background(), tt.product)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			tt.validate(t, saved)
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name        string
		products    []model.Product
		filterCat   enum.Category
		expectCount int
	}{
		{
			name: "success - only the products matching that category should be returned",
			products: []model.Product{
				{Name: "Cheeseburger", Category: enum.Meal},
				{Name: "Coca-Cola", Category: enum.Drink},
			},
			filterCat:   enum.Meal,
			expectCount: 1,
		},
		{
			name: "success - when no category is received, then all products should be returned regardless of their category",
			products: []model.Product{
				{Name: "Cheeseburger", Category: enum.Meal},
				{Name: "Coca-Cola", Category: enum.Drink},
				{Name: "Brownie", Category: enum.Dessert},
			},
			filterCat:   enum.Meal,
			expectCount: 3,
		},
		{
			name: "success - when has no products match the filter then should return an empty list",
			products: []model.Product{
				{Name: "Coca-Cola", Category: enum.Drink},
				{Name: "Pepsi", Category: enum.Drink},
			},
			filterCat:   enum.Meal,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := postgre.New(db)

			for _, p := range tt.products {
				require.NoError(t, db.Create(p.Build()).Error)
			}

			products, err := repo.GetAll(context.Background(), tt.filterCat)
			assert.NoError(t, err)
			assert.Len(t, products, tt.expectCount)
		})
	}
}

func TestUpdate(t *testing.T) {
	createProduct := func(t *testing.T, db *gorm.DB) model.Product {
		p := model.Product{
			Name:          "Cheeseburger",
			Price:         24.90,
			Description:   "Delicious cheeseburger",
			PreparingTime: 15,
			Category:      enum.Meal,
			ImageURL:      "https://example.com/images/cheeseburger.png",
		}
		require.NoError(t, db.Create(&p).Error)
		return p
	}

	tests := []struct {
		name      string
		setup     func(t *testing.T, db *gorm.DB) string
		update    model.Product
		expectErr bool
		validate  func(t *testing.T, updated model.Product, original model.Product)
	}{
		{
			name: "success - update some fields",
			setup: func(t *testing.T, db *gorm.DB) string {
				p := createProduct(t, db)
				return p.ID
			},
			update: model.Product{Name: "Updated Name", Price: 30.0},
			validate: func(t *testing.T, updated model.Product, original model.Product) {
				assert.Equal(t, "Updated Name", updated.Name)
				assert.Equal(t, 30.0, updated.Price)
				assert.Equal(t, original.Description, updated.Description)
			},
		},
		{
			name: "error - product not found",
			setup: func(t *testing.T, db *gorm.DB) string {
				return "non-existent-id"
			},
			update:    model.Product{Name: "Whatever"},
			expectErr: true,
		},
		{
			name: "no updates - returns existing product",
			setup: func(t *testing.T, db *gorm.DB) string {
				p := createProduct(t, db)
				return p.ID
			},
			update: model.Product{},
			validate: func(t *testing.T, updated model.Product, original model.Product) {
				assert.Equal(t, original.Name, updated.Name)
				assert.Equal(t, original.ID, updated.ID)
			},
		},
		{
			name: "error - db update fails",
			setup: func(t *testing.T, db *gorm.DB) string {
				return "invalid-uuid-format"
			},
			update:    model.Product{Name: "Name"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := postgre.New(db)

			original := model.Product{}
			id := tt.setup(t, db)
			if id != "non-existent-id" && id != "invalid-uuid-format" {
				_ = db.First(&original, "id = ?", id)
			}

			updated, err := repo.Update(context.Background(), id, tt.update)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				tt.validate(t, updated, original)
			}
		})
	}
}

func TestFindByID(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, db *gorm.DB) string
		expectErr  bool
		expectedID string
	}{
		{
			name: "success - product found",
			setup: func(t *testing.T, db *gorm.DB) string {
				p := model.Product{
					Name:          "Cheeseburger",
					Price:         24.90,
					Description:   "Delicious cheeseburger",
					PreparingTime: 15,
					Category:      enum.Meal,
					ImageURL:      "https://example.com/images/cheeseburger.png",
				}.Build()
				require.NoError(t, db.Create(&p).Error)
				return p.ID
			},
			expectErr: false,
		},
		{
			name: "error - product not found",
			setup: func(t *testing.T, db *gorm.DB) string {
				return "non-existent-id"
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := postgre.New(db)

			id := tt.setup(t, db)

			found, err := repo.FindByID(context.Background(), id)

			if tt.expectErr {
				require.Error(t, err)
				var notFoundErr *apperror.NotFoundError
				assert.True(t, errors.As(err, &notFoundErr))
				assert.Equal(t, "Product not found", notFoundErr.Msg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, id, found.ID)
			}
		})
	}
}

func TestFindByIDs(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, db *gorm.DB) []string
		expectErr  bool
		expectSize int
	}{
		{
			name: "success - products found",
			setup: func(t *testing.T, db *gorm.DB) []string {
				p1 := model.Product{Name: "Cheeseburger", Category: enum.Meal}.Build()
				p2 := model.Product{Name: "Coca-Cola", Category: enum.Drink}.Build()
				require.NoError(t, db.Create(&p1).Error)
				require.NoError(t, db.Create(&p2).Error)
				return []string{p1.ID, p2.ID}
			},
			expectErr:  false,
			expectSize: 2,
		},
		{
			name: "error - no products found",
			setup: func(t *testing.T, db *gorm.DB) []string {
				return []string{"non-existent-id1", "non-existent-id2"}
			},
			expectErr:  true,
			expectSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := postgre.New(db)

			ids := tt.setup(t, db)

			products, err := repo.FindByIDs(context.Background(), ids)

			if tt.expectErr {
				require.Error(t, err)
				assert.Nil(t, products)
				var notFoundErr *apperror.NotFoundError
				assert.True(t, errors.As(err, &notFoundErr))
				assert.Equal(t, "No products found", notFoundErr.Msg)
			} else {
				require.NoError(t, err)
				assert.Len(t, products, tt.expectSize)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T, db *gorm.DB) string
		expectErr bool
	}{
		{
			name: "success - delete existing product",
			setup: func(t *testing.T, db *gorm.DB) string {
				p := model.Product{Name: "To be deleted", Category: enum.Meal}.Build()
				require.NoError(t, db.Create(&p).Error)
				return p.ID
			},
			expectErr: false,
		},
		{
			name: "error - product not found",
			setup: func(t *testing.T, db *gorm.DB) string {
				return "non-existent-id"
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := postgre.New(db)

			id := tt.setup(t, db)

			err := repo.Delete(context.Background(), id)

			if tt.expectErr {
				assert.Error(t, err)
				var notFoundErr *apperror.NotFoundError
				assert.True(t, errors.As(err, &notFoundErr))
				assert.Equal(t, "Product not found", notFoundErr.Msg)
			} else {
				assert.NoError(t, err)
				var count int64
				db.Model(&model.Product{}).Where("id = ?", id).Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}
