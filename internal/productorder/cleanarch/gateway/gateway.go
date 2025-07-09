package gateway

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/external/datasource"
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

func (g *Gateway) CreateBulk(c context.Context, productOrders []entity.ProductOrder) (int, error) {
	var listProductOrderDAO = entity.ToListProducOrderDAO(productOrders)
	length, err := g.Datasource.CreateBulk(c, listProductOrderDAO)

	if err != nil {
		return 0, &apperror.InternalError{Msg: err.Error()}
	}

	return length, nil
}
