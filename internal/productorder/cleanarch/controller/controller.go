package controller

import "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/external/datasource"

type Controller struct {
	ProductOrderDatasource datasource.DataSource
}

func Build(productDataSource datasource.DataSource) *Controller {
	return &Controller{
		ProductOrderDatasource: productDataSource}
}
