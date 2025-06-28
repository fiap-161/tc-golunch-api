package usecases

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/gateway"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type UseCases struct {
	ProductGateway gateway.Gateway
}

func Build(productGateway gateway.Gateway) *UseCases {
	return &UseCases{ProductGateway: productGateway}
}

func (u *UseCases) CreateProduct(ctx context.Context, productDTO dto.ProductRequestDTO) (entity.Product, error) {
	var product entity.Product
	product = product.FromRequestDTO(productDTO)
	isValidCategory := enum.IsValidCategory(string(product.Category))

	if !isValidCategory {
		return entity.Product{}, &apperror.ValidationError{Msg: "Invalid category"}
	}

	saved, err := u.ProductGateway.Create(ctx, product.Build())
	if err != nil {
		return entity.Product{}, &apperror.InternalError{Msg: err.Error()}
	}

	return saved, nil
}

func (u *UseCases) ListCategories(ctx context.Context) []enum.Category {
	return enum.GetAllCategories()
}

func (u *UseCases) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	uploadDir := os.Getenv("UPLOAD_DIR")
	publicURL := os.Getenv("PUBLIC_URL")

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}
	contentType := http.DetectContentType(buffer)

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	if !allowedTypes[contentType] {
		return "", &apperror.ValidationError{Msg: "Only JPEG and PNG images are allowed"}
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("error resetting file: %w", err)
	}

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return "", fmt.Errorf("error creating upload dir: %w", err)
		}
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
	fullPath := filepath.Join(uploadDir, filename)

	if err := saveFile(fileHeader, fullPath); err != nil {
		return "", fmt.Errorf("error saving file: %w", err)
	}

	return fmt.Sprintf("%s/uploads/%s", publicURL, filename), nil
}

func saveFile(fileHeader *multipart.FileHeader, dest string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
