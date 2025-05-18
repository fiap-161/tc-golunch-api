package main

import (
	"log"
	"os"

	_ "github.com/fiap-161/tech-challenge-fiap161/docs"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivens/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivens/postgre"
	restProduct "github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/rest"
	servicesProduct "github.com/fiap-161/tech-challenge-fiap161/internal/product/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           GoLunch
// @version         1.0
// @description     Rest API para facilitar o gerenciamento de pedidos em uma lanchonete
// @host            localhost:8080
// @BasePath        /
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	r := gin.Default()

	DB, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	productRepository := postgre.NewProductRepository(DB)
	productService := servicesProduct.NewProductService(productRepository)
	productHandler := restProduct.NewProductHandler(productService)

	DB.AutoMigrate(&dto.ProductDAO{})

	// registering api routes
	r.POST("/product", productHandler.Create)
	r.GET("/product/categories", productHandler.ListCategories)
	r.GET("/product", productHandler.GetAll)
	r.PUT("/product/:id", productHandler.Update)
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
