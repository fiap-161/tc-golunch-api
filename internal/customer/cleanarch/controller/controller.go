package controller

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/cleanarch/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/cleanarch/presenter"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/cleanarch/usecases"
)

type Controller struct {
	CustomerDataSource datasource.DataSource
}

func Build(Customer datasource.DataSource) *Controller {
	return &Controller{
		CustomerDataSource: customerDataSource}
}
