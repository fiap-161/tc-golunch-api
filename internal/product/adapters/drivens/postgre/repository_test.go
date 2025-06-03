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
	t.Run("success - inserts valid product without error", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		product := model.Product{
			Name:          "Cheeseburger",
			Price:         24.90,
			Description:   "Delicioso cheeseburger com duas carnes.",
			PreparingTime: 15,
			Category:      enum.Category(enum.Meal),
			ImageURL:      "https://example.com/images/cheeseburger.png",
		}.Build()

		saved, err := repo.Create(context.Background(), product)

		assert.NoError(t, err)
		assert.NotEmpty(t, saved.ID)
		assert.Equal(t, product.Name, saved.Name)
	})

	t.Run("error - missing required fields", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		invalidProduct := model.Product{}

		saved, err := repo.Create(context.Background(), invalidProduct)

		assert.Error(t, err)
		assert.EqualError(t, err, "name is required")
		assert.Empty(t, saved.ID)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("success - only the products matching that category should be returned", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		db.Create(model.Product{
			Name:     "Cheeseburger",
			Category: enum.Meal,
		}.Build())
		db.Create(model.Product{
			Name:     "Coca-Cola",
			Category: enum.Drink,
		}.Build())

		products, err := repo.GetAll(context.Background(), uint(enum.Meal))

		assert.NoError(t, err)
		assert.Len(t, products, 1)
		assert.Equal(t, enum.Meal, products[0].Category)
	})

	t.Run("success - when no category is received, then all products should be returned regardless of their category", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		db.Create(model.Product{
			Name:     "Cheeseburger",
			Category: enum.Meal,
		}.Build())
		db.Create(model.Product{
			Name:     "Coca-Cola",
			Category: enum.Drink,
		}.Build())
		db.Create(model.Product{
			Name:     "Brownie",
			Category: enum.Dessert,
		}.Build())

		products, err := repo.GetAll(context.Background(), 0)

		assert.NoError(t, err)
		assert.Len(t, products, 3)
	})

	t.Run("success - when has no products match the filter then should return an empty list", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		db.Create(model.Product{
			Name:     "Coca-Cola",
			Category: enum.Drink,
		}.Build())
		db.Create(model.Product{
			Name:     "Pepsi",
			Category: enum.Drink,
		}.Build())

		products, err := repo.GetAll(context.Background(), uint(enum.Meal))

		assert.NoError(t, err)
		assert.Len(t, products, 0)
	})
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

	t.Run("success - update some fields", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		original := createProduct(t, db)

		update := model.Product{
			Name:  "Updated Name",
			Price: 30.0,
		}

		updated, err := repo.Update(context.Background(), original.ID, update)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)
		assert.Equal(t, 30.0, updated.Price)
		assert.Equal(t, original.Description, updated.Description)
	})

	t.Run("error - product not found", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		_, err := repo.Update(context.Background(), "non-existent-id", model.Product{Name: "Whatever"})
		assert.Error(t, err)
	})

	t.Run("no updates - returns existing product", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		original := createProduct(t, db)

		emptyUpdate := model.Product{} // no fields to update
		existing, err := repo.Update(context.Background(), original.ID, emptyUpdate)
		require.NoError(t, err)
		assert.Equal(t, original.Name, existing.Name)
		assert.Equal(t, original.ID, existing.ID)
	})

	t.Run("error - db update fails", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		badID := "invalid-uuid-format"
		update := model.Product{Name: "Name"}

		_, err := repo.Update(context.Background(), badID, update)
		assert.Error(t, err)
	})
}

func TestFindByID(t *testing.T) {
	t.Run("success - product found", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		product := model.Product{
			Name:          "Cheeseburger",
			Price:         24.90,
			Description:   "Delicious cheeseburger",
			PreparingTime: 15,
			Category:      enum.Meal,
			ImageURL:      "https://example.com/images/cheeseburger.png",
		}
		err := db.Create(&product).Error
		require.NoError(t, err)

		found, err := repo.FindByID(context.Background(), product.ID)
		require.NoError(t, err)
		assert.Equal(t, product.ID, found.ID)
		assert.Equal(t, product.Name, found.Name)
	})

	t.Run("error - product not found", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		_, err := repo.FindByID(context.Background(), "non-existent-id")
		require.Error(t, err)

		var notFoundErr *apperror.NotFoundError
		assert.True(t, errors.As(err, &notFoundErr))
		assert.Equal(t, "Product not found", notFoundErr.Msg)
	})

}

func TestFindByIDs(t *testing.T) {
	t.Run("success - products found", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		p1 := model.Product{
			Name:     "Cheeseburger",
			Category: enum.Category(enum.Meal)}.Build()

		p2 := model.Product{
			Name:     "Coca-cola",
			Category: enum.Drink,
		}.Build()

		require.NoError(t, db.Create(&p1).Error)
		require.NoError(t, db.Create(&p2).Error)

		ids := []string{p1.ID, p2.ID}

		products, err := repo.FindByIDs(context.Background(), ids)
		require.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Contains(t, []string{products[0].ID, products[1].ID}, p1.ID)
		assert.Contains(t, []string{products[0].ID, products[1].ID}, p2.ID)
	})

	t.Run("error - no products found", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		ids := []string{"non-existent-id1", "non-existent-id2"}

		products, err := repo.FindByIDs(context.Background(), ids)
		require.Error(t, err)
		assert.Nil(t, products)

		var notFoundErr *apperror.NotFoundError
		assert.True(t, errors.As(err, &notFoundErr))
		assert.Equal(t, "No products found", notFoundErr.Msg)
	})
}

func TestDelete(t *testing.T) {
	t.Run("success - delete existing product", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		product := model.Product{
			Name:     "To be deleted",
			Category: enum.Meal,
		}.Build()
		require.NoError(t, db.Create(&product).Error)

		err := repo.Delete(context.Background(), product.ID)
		assert.NoError(t, err)

		// Verificar se realmente deletou
		var count int64
		db.Model(&model.Product{}).Where("id = ?", product.ID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("error - product not found", func(t *testing.T) {
		db := setupTestDB(t)
		repo := postgre.New(db)

		err := repo.Delete(context.Background(), "non-existent-id")
		assert.Error(t, err)

		var notFoundErr *apperror.NotFoundError
		assert.True(t, errors.As(err, &notFoundErr))
		assert.Equal(t, "Product not found", notFoundErr.Msg)
	})
}
