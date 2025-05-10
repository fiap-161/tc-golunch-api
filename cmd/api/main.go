package main

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/user/adapters/drivens/postgre"
	"github.com/fiap-161/tech-challenge-fiap161/internal/user/adapters/drivers/rest"
	"github.com/fiap-161/tech-challenge-fiap161/internal/user/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/user/service"
	"log"
	"os"

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

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("error to migrate: %v", err)
	}

	userRepository := postgre.NewRepository(db)
	userService := service.New(userRepository)
	userHandler := rest.NewUserHandler(userService)

	// registering api routes
	r.GET("/user/:id", userHandler.GetUserByID)
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
