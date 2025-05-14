package main

import (
	_ "github.com/fiap-161/tech-challenge-fiap161/docs"
	restProduct "github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/rest"
	servicesProduct "github.com/fiap-161/tech-challenge-fiap161/internal/product/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           GoLunch
// @version         1.0
// @description     Rest API para facilitar o gerenciamento de pedidos em uma lanchonete
// @host            localhost:8080
// @BasePath        /
func main() {
	r := gin.Default()

	/*
		_, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	*/

	productService := servicesProduct.NewProductService()
	productHandler := restProduct.NewProductHandler(productService)

	// registering api routes
	r.POST("/product", productHandler.Create)
	r.GET("/product/categories", productHandler.ListCategories)
	r.GET("/ping", ping)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}

// Ping godoc
// @Summary      Responde com "Pong"
// @Description  Health Check
// @Tags         Ping
// @Accept       json
// @Produce      json
// @Success      200 {object}  PongResponse
// @Router       /ping [get]
func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type PongResponse struct {
	Message string `json:"message" example:"pong"`
}
