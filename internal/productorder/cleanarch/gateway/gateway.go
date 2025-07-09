package gateway

import "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/external/datasource"

type Gateway struct {
	Datasource datasource.DataSource
}

func Build(datasource datasource.DataSource) *Gateway {
	return &Gateway{
		Datasource: datasource,
	}
}
