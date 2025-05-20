package main

import (
	"log"
	"os"
	"time"

	adminPostgre "github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivens/postgre"
	adminRest "github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivers/rest"
	admin "github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/model"
	adminService "github.com/fiap-161/tech-challenge-fiap161/internal/admin/service"
	auth "github.com/fiap-161/tech-challenge-fiap161/internal/auth/adapters/jwt"
	customerPostgre "github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivens/postgre"
	customerRest "github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest"
	customer "github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
	customerService "github.com/fiap-161/tech-challenge-fiap161/internal/customer/service"
	"github.com/fiap-161/tech-challenge-fiap161/internal/http/middleware"

	_ "github.com/fiap-161/tech-challenge-fiap161/docs"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivens/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivens/postgre"
	restProduct "github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/rest"
	servicesProduct "github.com/fiap-161/tech-challenge-fiap161/internal/product/services"
	"github.com/gin-gonic/gin"
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// DESCOMENTAR PARA RODAR APENAS O BANCO NO DOCKER
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Erro ao carregar o .env")
	// }

	r := gin.Default()

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&customer.Customer{}, &admin.Admin{}, &dto.Product{})
	if err != nil {
		log.Fatalf("error to migrate: %v", err)
	}

	// jwt service for generate and validate tokens
	jwtService := auth.NewJWTService(os.Getenv("SECRET_KEY"), 24*time.Hour)

	// customer
	customerRepository := customerPostgre.NewRepository(db)
	customerService := customerService.New(customerRepository, jwtService)
	customerHandler := customerRest.NewCustomerHandler(customerService)

	//admin
	adminRepository := adminPostgre.NewRepository(db)
	adminService := adminService.New(adminRepository, jwtService)
	adminHandler := adminRest.NewAdminHandler(adminService)

	//product
	productRepository := postgre.NewProductRepository(db)
	productService := servicesProduct.NewProductService(productRepository)
	productHandler := restProduct.NewProductHandler(productService)

	// Rotas default
	r.GET("/ping", ping)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rotas públicas (login/register)
	r.GET("/identify/:cpf", customerHandler.Identify)
	r.GET("/anonymous", customerHandler.Anonymous)
	r.POST("/customer/register", customerHandler.Create)
	r.POST("/admin/register", adminHandler.Register)
	r.POST("/admin/login", adminHandler.Login)

	// Grupo autenticado
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware(jwtService))

	// Rotas acessíveis para qualquer usuário autenticado
	// Produto
	authenticated.GET("/product/categories", productHandler.ListCategories)
	authenticated.GET("/product", productHandler.GetAll)

	// Grupo para admins dentro do grupo autenticado
	adminRoutes := authenticated.Group("/product")
	adminRoutes.Use(middleware.AdminOnly())
	adminRoutes.POST("/", productHandler.Create)
	adminRoutes.PUT("/:id", productHandler.ValidateIfProductExists, productHandler.Update)
	adminRoutes.DELETE("/:id", productHandler.ValidateIfProductExists, productHandler.Delete)

	r.Run(":8080")
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
