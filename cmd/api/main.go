package main

import (
	adminPostgre "github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivens/postgre"
	adminRest "github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivers/rest"
	admin "github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/model"
	adminService "github.com/fiap-161/tech-challenge-fiap161/internal/admin/service"
	"github.com/fiap-161/tech-challenge-fiap161/internal/auth"
	customerPostgre "github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivens/postgre"
	customerRest "github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest"
	customer "github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
	customerService "github.com/fiap-161/tech-challenge-fiap161/internal/customer/service"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/fiap-161/tech-challenge-fiap161/docs"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

// @title           GoLunch
// @version         1.0
// @description     Rest API para facilitar o gerenciamento de pedidos em uma lanchonete
// @host            localhost:8080
// @BasePath        /
func main() {
	r := gin.Default()

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&customer.Customer{}, &admin.Admin{})
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

	// customer routes
	r.GET("/identify/:cpf", customerHandler.Identify)
	r.GET("/anonymous", customerHandler.Anonymous)
	r.POST("/customer/register", customerHandler.Create)

	//admin routes
	r.POST("/admin/register", adminHandler.Register)
	r.POST("/admin/login", adminHandler.Login)

	//api default routes
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
