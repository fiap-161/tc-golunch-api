package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/fiap-161/tech-challenge-fiap161/database"
	_ "github.com/fiap-161/tech-challenge-fiap161/docs"

	// Product Clean Architecture
	productController "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/controller"
	productDataSource "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/external/datasource"
	productFactory "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/factory"
	productHandler "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/handler"

	// Order Clean Architecture
	orderController "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/controller"
	orderDataSource "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/external/datasource"
	orderFactory "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/factory"
	orderHandler "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/handler"

	// Payment Clean Architecture
	paymentController "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/controllers"
	paymentDataSource "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/external/datasource"
	paymentFactory "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/factory"
	paymentHandler "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/handlers"

	// Auth
	authController "github.com/fiap-161/tech-challenge-fiap161/internal/auth/cleanarch/controller"
	"github.com/fiap-161/tech-challenge-fiap161/internal/auth/cleanarch/external"
	"github.com/fiap-161/tech-challenge-fiap161/internal/http/middleware"
)

// @title           GoLunch
// @version         1.0
// @description     REST API to facilitate order management in a snack bar.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := gin.Default()
	loadYAML()

	db := database.NewPostgresDatabase().GetDb()

	// Auto migrate database
	if err := db.AutoMigrate(
	// ... migrations
	); err != nil {
		log.Fatalf("Erro ao migrar o banco: %v", err)
	}

	// JWT service for generate and validate tokens
	jwtGateway := external.NewJWTService(os.Getenv("SECRET_KEY"), 24*time.Hour)
	authController := authController.New(jwtGateway)

	// Build Product Domain
	productDataSource := productDataSource.New(db)
	productUseCases := productFactory.BuildProductUseCases(productDataSource)
	productController := productController.Build(productUseCases)
	productHandler := productHandler.New(productController)

	// Build Order Domain
	orderDataSource := orderDataSource.New(db)
	orderUseCases := orderFactory.BuildOrderUseCases(orderDataSource, productUseCases, nil, nil) // Dependencies injected
	orderController := orderController.Build(orderUseCases)
	orderHandler := orderHandler.New(orderController)

	// Build Payment Domain
	paymentDataSource := paymentDataSource.New(db)
	paymentUseCases := paymentFactory.BuildPaymentUseCases(paymentDataSource, nil, nil, nil, nil) // Dependencies injected
	paymentController := paymentController.Build(paymentUseCases)
	paymentHandler := paymentHandler.New(paymentController)

	// Setup routes
	setupRoutes(r, productHandler, orderHandler, paymentHandler, authController)

	r.Run(":8080")
}

func setupRoutes(r *gin.Engine, productHandler, orderHandler, paymentHandler interface{}, authController interface{}) {
	// Setup routes here
	// This is just an example structure
}

func loadYAML() {
	viper.SetConfigName("default")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./conf/environment")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo de configuração:", err)
	}
}
