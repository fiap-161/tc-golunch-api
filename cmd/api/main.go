package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/fiap-161/tech-challenge-fiap161/database"
	_ "github.com/fiap-161/tech-challenge-fiap161/docs"
	adminpostgre "github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivens/postgre"
	adminrest "github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivers/rest"
	adminmodel "github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/model"
	adminservice "github.com/fiap-161/tech-challenge-fiap161/internal/admin/service"
	auth "github.com/fiap-161/tech-challenge-fiap161/internal/auth/adapters/jwt"
	customerpostgre "github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivens/postgre"
	customerrest "github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest"
	customermodel "github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
	customerservice "github.com/fiap-161/tech-challenge-fiap161/internal/customer/service"
	"github.com/fiap-161/tech-challenge-fiap161/internal/http/middleware"
	orderpostgre "github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivens/postgre"
	orderrest "github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivers/rest"
	order "github.com/fiap-161/tech-challenge-fiap161/internal/order/core/model"
	orderservice "github.com/fiap-161/tech-challenge-fiap161/internal/order/service"
	paymentpostgre "github.com/fiap-161/tech-challenge-fiap161/internal/payment/adapters/drivens/postgre"
	paymenthandler "github.com/fiap-161/tech-challenge-fiap161/internal/payment/adapters/drivers/rest"
	paymentmodel "github.com/fiap-161/tech-challenge-fiap161/internal/payment/core/model"
	paymentservice "github.com/fiap-161/tech-challenge-fiap161/internal/payment/service"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/controller"
	productmodel "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/handler"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/presenter"
	productpostgre "github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/adapters/drivens/postgre"
	restproduct "github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/adapters/drivers/rest"
	servicesproduct "github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/service"
	productorderrepository "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/adapters/drivens/postgre"
	productordermodel "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/adapters/mercadopago"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	r := gin.Default()
	loadYAML()

	db := database.NewPostgresDatabase().GetDb()

	if err := db.AutoMigrate(
		&customermodel.Customer{},
		&adminmodel.Admin{},
		&productmodel.ProductDAO{},
		&order.Order{},
		&productordermodel.ProductOrder{},
		&paymentmodel.Payment{},
	); err != nil {
		log.Fatalf("Erro ao migrar o banco: %v", err)
	}

	// servir arquivos estáticos - imagens
	uploadDir := os.Getenv("UPLOAD_DIR")

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
	productRepository := productpostgre.New(db)
	productService := servicesproduct.New(productRepository)
	productHandler := restproduct.New(productService)

	// ProductOrder
	productOrderRepository := productorderrepository.New(db)

	// QR Code Client
	qrCodeClient := mercadopago.New()

	// Order Repository
	orderRepository := orderpostgre.New(db)

	// Payment
	paymentRepository := paymentpostgre.New(db)
	paymentService := paymentservice.New(
		qrCodeClient,
		orderRepository,
		paymentRepository,
		productOrderRepository,
		productRepository,
	)
	paymentHandler := paymenthandler.New(paymentService)

	// Order Service
	orderService := orderservice.New(
		orderRepository,
		productRepository,
		productOrderRepository,
		paymentService,
	)
	orderHandler := orderrest.New(orderService)

	// CLEAN ARCH - Product
	productDataSource := datasource.New(db)
	productPresenter := presenter.Build()
	productController := controller.Build(productDataSource, *productPresenter)
	productHandlerCleanArch := handler.New(productController)

	// Default Routes
	r.GET("/ping", ping)
	r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	r.Static("/uploads", uploadDir)

	// Public Routes (login/register)
	r.GET("/customer/identify/:cpf", customerHandler.Identify)
	r.GET("/customer/anonymous", customerHandler.Anonymous)
	r.POST("/customer/register", customerHandler.Create)
	r.POST("/admin/register", adminHandler.Register)
	r.POST("/admin/login", adminHandler.Login)

	// Webhook for Mercado Pago
	r.POST("/webhook/payment/check", paymentHandler.CheckPayment)

	// Authenticated Group
	authenticated := r.Group("/")
	authenticated.Use(middleware.AuthMiddleware(jwtService))

	// Routes for regular authenticated users
	// Product
	authenticated.GET("/product/categories", productHandlerCleanArch.ListCategories)
	authenticated.GET("/product", productHandlerCleanArch.GetAllByCategory)

	// Order
	authenticated.POST("/order", orderHandler.Create)
	authenticated.GET("/order", middleware.AdminOnly(), orderHandler.GetAll)
	authenticated.PUT("/order/:id", middleware.AdminOnly(), orderHandler.Update)
	authenticated.GET("/order/panel", middleware.AdminOnly(), orderHandler.GetPanel)

	// Group for admin users inside authenticated group
	adminRoutes := authenticated.Group("/product")
	adminRoutes.Use(middleware.AdminOnly())
	adminRoutes.POST("/image/upload", productHandlerCleanArch.UploadImage)
	adminRoutes.POST("/", productHandlerCleanArch.Create)
	adminRoutes.PUT("/:id", productHandlerCleanArch.Update)
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

type PongResponse struct {
	Message string `json:"message"`
}
