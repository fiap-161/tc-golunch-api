package controller

import (
	"context"
	"mime/multipart"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/presenter"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/usecases"
)

// INFO: controllers Criam gateways e requisitam usecases
type Controller struct {
	ProductDatasource datasource.DataSource
}

func Build(productDataSource datasource.DataSource) *Controller {
	return &Controller{
		ProductDatasource: productDataSource}
}

func (c *Controller) Create(ctx context.Context, productDTO dto.ProductRequestDTO) (dto.ProductResponseDTO, error) {
	productGateway := gateway.Build(c.ProductDatasource)
	useCase := usecases.Build(*productGateway)
	presenter := presenter.Build()

	var product entity.Product
	product = entity.FromRequestDTO(productDTO)
	product, err := useCase.CreateProduct(ctx, product)

	if err != nil {
		return dto.ProductResponseDTO{}, err
	}

	return presenter.FromEntityToResponseDTO(product), nil
}

func (c *Controller) ListCategories(ctx context.Context) []enum.Category {
	productGateway := gateway.Build(c.ProductDatasource)
	useCase := usecases.Build(*productGateway)
	return useCase.ListCategories(ctx)
}

func (c *Controller) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	productGateway := gateway.Build(c.ProductDatasource)
	useCase := usecases.Build(*productGateway)
	return useCase.UploadImage(ctx, fileHeader)
}

func (c *Controller) GetAllByCategory(ctx context.Context, category string) (dto.ProductListResponseDTO, error) {
	productGateway := gateway.Build(c.ProductDatasource)
	useCase := usecases.Build(*productGateway)
	presenter := presenter.Build()

	result, err := useCase.GetAllByCategory(ctx, category)

	if err != nil {
		return dto.ProductListResponseDTO{}, err
	}

	return presenter.FromEntityListToProductListResponseDTO(result), nil
}

func (c *Controller) Update(ctx context.Context, productId string, productDTO dto.ProductRequestUpdateDTO) (dto.ProductResponseDTO, error) {
	productGateway := gateway.Build(c.ProductDatasource)
	useCase := usecases.Build(*productGateway)
	presenter := presenter.Build()

	product := entity.FromUpdateDTO(productDTO)
	result, err := useCase.Update(ctx, productId, product)

	if err != nil {
		return dto.ProductResponseDTO{}, err
	}

	return presenter.FromEntityToResponseDTO(result), nil
}

func (c *Controller) FindByID(ctx context.Context, productId string) (dto.ProductResponseDTO, error) {
	productGateway := gateway.Build(c.ProductDatasource)
	useCase := usecases.Build(*productGateway)
	presenter := presenter.Build()

	result, err := useCase.FindByID(ctx, productId)

	if err != nil {
		return dto.ProductResponseDTO{}, err
	}

	return presenter.FromEntityToResponseDTO(result), nil
}

func (c *Controller) Delete(ctx context.Context, productId string) error {
	productGateway := gateway.Build(c.ProductDatasource)
	useCase := usecases.Build(*productGateway)

	err := useCase.Delete(ctx, productId)

	if err != nil {
		return err
	}

	return nil
}
