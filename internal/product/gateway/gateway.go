package gateway

import (
	"context"
	"errors"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/external/datasource"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type Gateway struct {
	Datasource datasource.DataSource
}

func Build(datasource datasource.DataSource) *Gateway {
	return &Gateway{
		Datasource: datasource,
	}
}

func (g *Gateway) Create(c context.Context, product entity.Product) (entity.Product, error) {
	var productDAO = dto.ToProductDAO(product)
	created, err := g.Datasource.Create(c, productDAO)

	if err != nil {
		return entity.Product{}, &apperror.InternalError{Msg: err.Error()}
	}

	return dto.FromProductDAO(created), nil
}

func (g *Gateway) GetAllByCategory(c context.Context, category string) ([]entity.Product, error) {
	result, err := g.Datasource.GetAllByCategory(c, category)

	if err != nil {
		return []entity.Product{}, &apperror.InternalError{Msg: err.Error()}
	}

	var products []entity.Product
	for _, dao := range result {
		entity := dto.FromProductDAO(dao)
		products = append(products, entity)
	}

	return products, nil
}

func (g *Gateway) Update(c context.Context, productId string, product entity.Product) (entity.Product, error) {
	productDAO := dto.ToProductDAO(product)
	updated, err := g.Datasource.Update(c, productId, productDAO)

	if err != nil {
		return entity.Product{}, &apperror.InternalError{Msg: err.Error()}
	}

	return dto.FromProductDAO(updated), nil
}

func (g *Gateway) FindByID(c context.Context, productId string) (entity.Product, error) {
	found, err := g.Datasource.FindByID(c, productId)

	if err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			return entity.Product{}, notFoundErr
		}
		return entity.Product{}, &apperror.InternalError{Msg: "Unexpected error"}
	}

	return dto.FromProductDAO(found), nil
}

func (g *Gateway) Delete(c context.Context, productId string) error {
	err := g.Datasource.Delete(c, productId)

	if err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			return notFoundErr
		}
		return &apperror.InternalError{Msg: "Unexpected error"}
	}

	return nil
}

func (g *Gateway) FindByIDs(c context.Context, productIdList []string) ([]entity.Product, error) {
	foundList, err := g.Datasource.FindByIDs(c, productIdList)

	if err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			return []entity.Product{}, notFoundErr
		}
		return []entity.Product{}, &apperror.InternalError{Msg: "Unexpected error"}
	}

	return dto.EntityListFromDAOList(foundList), nil
}
