package postgre_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"

	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/adapters/drivens/postgre"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/core/model"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Payment{})
	require.NoError(t, err)

	return db
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		payment  model.Payment
		wantErr  string
		validate func(t *testing.T, saved model.Payment)
	}{
		{
			name:    "success - inserts valid payment without error",
			payment: model.Payment{}.Build("1931e833-d29e-43f7-8799-420c8dbcc89e", "qrcode"),
			wantErr: "",
			validate: func(t *testing.T, saved model.Payment) {
				assert.NotEmpty(t, saved.ID)
				assert.Equal(t, "1931e833-d29e-43f7-8799-420c8dbcc89e", saved.OrderID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := postgre.New(db)

			saved, err := repo.Create(context.Background(), tt.payment)

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

func TestFindByOrderID2(t *testing.T) {
	tests := []struct {
		name     string
		orderID  string
		setup    func(db *gorm.DB)
		wantErr  string
		validate func(t *testing.T, payment model.Payment)
	}{
		{
			name:    "success - payment found by OrderID",
			orderID: "1931e833-d29e-43f7-8799-420c8dbcc89e", // OrderID existente
			setup: func(db *gorm.DB) {
				payment := model.Payment{}.Build("1931e833-d29e-43f7-8799-420c8dbcc89e", "qrcode")
				require.NoError(t, db.Create(&payment).Error)
			},
			wantErr: "",
			validate: func(t *testing.T, payment model.Payment) {
				assert.NotEmpty(t, payment.ID)
				assert.Equal(t, "1931e833-d29e-43f7-8799-420c8dbcc89e", payment.OrderID)
				assert.Equal(t, "qrcode", payment.QrCode)
			},
		},
		{
			name:    "failure - payment not found by OrderID",
			orderID: "uuid-not-found",
			setup:   func(db *gorm.DB) {},
			wantErr: "record not found",
			validate: func(t *testing.T, payment model.Payment) {
				assert.Empty(t, payment.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)

			tt.setup(db)

			repo := postgre.New(db)

			payment, err := repo.FindByOrderID(context.Background(), tt.orderID)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			tt.validate(t, payment)
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name     string
		payment  model.Payment
		update   model.Payment
		wantErr  string
		validate func(t *testing.T, updated model.Payment)
	}{
		{
			name:    "success - update payment",
			payment: model.Payment{}.Build("uuid-1", "qrcode-1"),
			update: model.Payment{
				Status: model.Approved,
			},
			wantErr: "",
			validate: func(t *testing.T, updated model.Payment) {
				assert.NotEmpty(t, updated.ID)
				assert.Equal(t, "uuid-1", updated.OrderID)
				assert.Equal(t, model.Approved, updated.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			repo := postgre.New(db)

			saved, err := repo.Create(context.Background(), tt.payment)
			assert.NoError(t, err)

			saved.Status = tt.update.Status
			updated, err := repo.Update(context.Background(), saved)

			if tt.wantErr != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			tt.validate(t, updated)
		})
	}
}
