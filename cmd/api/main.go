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
	admincontroller "github.com/fiap-161/tech-challenge-fiap161/internal/admin/cleanarch/controller"
	adminmodel "github.com/fiap-161/tech-challenge-fiap161/internal/admin/cleanarch/dto"
	admindatasource "github.com/fiap-161/tech-challenge-fiap161/internal/admin/cleanarch/external/datasource"
	adminhandler "github.com/fiap-161/tech-challenge-fiap161/internal/admin/cleanarch/handler"
	authcontroller "github.com/fiap-161/tech-challenge-fiap161/internal/auth/cleanarch/controller"
	"github.com/fiap-161/tech-challenge-fiap161/internal/auth/cleanarch/external"
	customerpostgre "github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivens/postgre"
	customerrest "github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest"
	customermodel "github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
	customerservice "github.com/fiap-161/tech-challenge-fiap161/internal/customer/service"
	"github.com/fiap-161/tech-challenge-fiap161/internal/http/middleware"
	orderadapters "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/adapters"
	ordercontroller "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/controller"
	ordermodel "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/dto"
	orderdatasource "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/external/datasource"
	ordergateway "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/gateway"
	orderhandler "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/handler"
	orderusecases "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/usecases"
	paymentadapters "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/adapters"
	paymentcontroller "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/controllers"
	paymentmodel "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/dto"
	paymentdatasource "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/external/datasource"
	paymentgateway "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/gateway"
	paymenthandler "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/handlers"
	paymentusecases "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/usecases"
	productadapters "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/adapters"
	productcontroller "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/controller"
	productmodel "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/dto"
	productdatasource "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/external/datasource"
	productgateway "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/gateway"
	producthandler "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/handler"
	productusecases "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/usecases"
	productordermodel "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/dto"
	productorderdatasource "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/external/datasource"
	productordergateway "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/gateway"
	productorderusecases "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/usecases"
	qrcodeprovider "github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/gateways"
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

	db := database.NewPostgresDatabase().GetDb()

	if err := db.AutoMigrate(
		&customermodel.Customer{}, // todo
		&adminmodel.AdminDAO{},
		&productmodel.ProductDAO{},
		&ordermodel.OrderDAO{},
		&productordermodel.ProductOrderDAO{},
		&paymentmodel.PaymentDAO{},
	); err != nil {
		log.Fatalf("Erro ao migrar o banco: %v", err)
	}

	// Serve static files
	uploadDir := os.Getenv("UPLOAD_DIR")

	// Jwt service for generate and validate tokens
	jwtGateway := external.NewJWTService(os.Getenv("SECRET_KEY"), 24*time.Hour)
	authController := authcontroller.New(jwtGateway)

	// ADMIN
	adminDatasource := admindatasource.New(db)
	adminController := admincontroller.Build(adminDatasource)
	adminHandler := adminhandler.New(adminController, authController)

	// Customer
	customerRepository := customerpostgre.NewRepository(db)
	customerSrv := customerservice.New(customerRepository, authController)
	customerHandler := customerrest.NewCustomerHandler(customerSrv)

	// Product
	productDataSource := productdatasource.New(db)
	productController := productcontroller.Build(productDataSource)
	productHandlerCleanArch := producthandler.New(productController)

	// CLEAN ARCH ProductOrder
	productOrderDataSource := productorderdatasource.New(db)
	productOrderGateway := productordergateway.Build(productOrderDataSource)
	productOrderUseCase := productorderusecases.Build(*productOrderGateway)

	// CLEAN ARCH Payment
	paymentDataSource := paymentdatasource.New(db)
	paymentGateway := paymentgateway.Build(paymentDataSource)

	// QR Code Client
	qrCodeClient := qrcodeprovider.New()

	// CLEAN ARCH Order
	orderDataSource := orderdatasource.New(db)
	orderGateway := ordergateway.Build(orderDataSource)

	// Common Gateways
	productGateway := productgateway.Build(productDataSource)
	productUseCase := productusecases.Build(*productGateway)

	productServiceAdapter := productadapters.NewProductServiceAdapter(productUseCase)
	productOrderServiceAdapterForOrder, productOrderServiceAdapterForPayment := productordergateway.NewProductOrderServiceAdapter(productOrderUseCase)

	// Creating payment use case without orderService (to avoid circular dependency)
	paymentUseCaseWithoutOrder := paymentusecases.Build(paymentGateway, qrCodeClient, productServiceAdapter, productOrderServiceAdapterForPayment, nil)
	paymentServiceAdapter := paymentadapters.NewPaymentServiceAdapter(paymentUseCaseWithoutOrder)

	// Creating orderUseCase with productService and productOrderService (to avoid circular dependency)
	orderUseCase := orderusecases.Build(orderGateway, productServiceAdapter, productOrderServiceAdapterForOrder, paymentServiceAdapter)

	// Creating orderServiceAdapter with orderUseCase
	orderServiceAdapter := orderadapters.NewOrderServiceAdapter(orderUseCase)

	// Creating payment use case with orderServiceAdapter
	paymentUseCase := paymentusecases.Build(paymentGateway, qrCodeClient, productServiceAdapter, productOrderServiceAdapterForPayment, orderServiceAdapter)

	// Order
	orderController := ordercontroller.Build(orderUseCase)
	orderHandler := orderhandler.New(orderController)

	// Payment
	paymentController := paymentcontroller.Build(paymentUseCase)
	paymentHandler := paymenthandler.New(paymentController)

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
	authenticated.Use(middleware.AuthMiddleware(authController))

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
	adminRoutes.DELETE("/:id", productHandlerCleanArch.Delete)

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
