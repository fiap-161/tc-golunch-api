package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/fiap-161/tech-challenge-fiap161/docs"
	adminpostgre "github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivens/postgre"
	adminrest "github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivers/rest"
	admin "github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/model"
	adminservice "github.com/fiap-161/tech-challenge-fiap161/internal/admin/service"
	auth "github.com/fiap-161/tech-challenge-fiap161/internal/auth/adapters/jwt"
	customerpostgre "github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivens/postgre"
	customerrest "github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest"
	customer "github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
	customerservice "github.com/fiap-161/tech-challenge-fiap161/internal/customer/service"
	"github.com/fiap-161/tech-challenge-fiap161/internal/http/middleware"
	orderpostgre "github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivens/postgre"
	orderrest "github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivers/rest"
	order "github.com/fiap-161/tech-challenge-fiap161/internal/order/core/model"
	orderservice "github.com/fiap-161/tech-challenge-fiap161/internal/order/service"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivens/dto"
	productpostgre "github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivens/postgre"
	productrest "github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/rest"
	productservice "github.com/fiap-161/tech-challenge-fiap161/internal/product/services"
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

	// UNCOMMENT TO RUN ONLY THE DATABASE IN DOCKER
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Erro ao carregar o .env")
	// }

	r := gin.Default()
	loadYAML()

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&customer.Customer{}, &admin.Admin{}, &dto.Product{}, &order.Order{})
	if err != nil {
		log.Fatalf("error to migrate: %v", err)
	}

	// Jwt service for generate and validate tokens
	jwtService := auth.NewJWTService(os.Getenv("SECRET_KEY"), 24*time.Hour)

	// Customer
	customerRepository := customerpostgre.NewRepository(db)
	customerSrv := customerservice.New(customerRepository, jwtService)
	customerHandler := customerrest.NewCustomerHandler(customerSrv)

	// Admin
	adminRepository := adminpostgre.NewRepository(db)
	adminSrv := adminservice.New(adminRepository, jwtService)
	adminHandler := adminrest.NewAdminHandler(adminSrv)

	// Product
	productRepository := productpostgre.NewProductRepository(db)
	productService := productservice.NewProductService(productRepository)
	productHandler := productrest.New(productService)

	// Order
	orderRepository := orderpostgre.NewRepository(db)
	orderService := orderservice.New(orderRepository, productRepository)
	orderHandler := orderrest.NewOrderHandler(orderService)

	// Default Routes
	r.GET("/ping", ping)
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	// Public Routes (login/register)
	r.GET("/identify/:cpf", customerHandler.Identify)
	r.GET("/anonymous", customerHandler.Anonymous)
	r.POST("/customer/register", customerHandler.Create)
	r.POST("/admin/register", adminHandler.Register)
	r.POST("/admin/login", adminHandler.Login)

	// Authenticated Group
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware(jwtService))

	// Routes for regular authenticated users
	// Product
	authenticated.GET("/product/categories", productHandler.ListCategories)
	authenticated.GET("/product", productHandler.GetAll)

	// Order
	authenticated.POST("/order", orderHandler.Create)
	authenticated.GET("/order", orderHandler.GetAll)

	// Group for admin users inside authenticated group
	adminRoutes := authenticated.Group("/product")
	adminRoutes.Use(middleware.AdminOnly())
	adminRoutes.POST("/", productHandler.Create)
	adminRoutes.PUT("/:id", productHandler.ValidateIfProductExists, productHandler.Update)
	adminRoutes.DELETE("/:id", productHandler.ValidateIfProductExists, productHandler.Delete)

	r.Run(":8080")
}

func loadYAML() {
	viper.SetConfigName("default")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./conf/environment")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading yaml config: %v", err)
	}
}

// Ping godoc
// @Summary      Answers with "pong"
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
