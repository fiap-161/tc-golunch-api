package database

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

func NewPostgresDatabase() Database {
	once.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			// Build DSN from individual environment variables for Kubernetes
			host := os.Getenv("DATABASE_HOST")
			port := os.Getenv("DATABASE_PORT")
			dbname := os.Getenv("DATABASE_NAME")
			user := os.Getenv("DATABASE_USER")
			password := os.Getenv("DATABASE_PASSWORD")
			
			if host == "" || port == "" || dbname == "" || user == "" || password == "" {
				log.Fatal("DATABASE environment variables not found. Please set DATABASE_URL or DATABASE_HOST, DATABASE_PORT, DATABASE_NAME, DATABASE_USER, DATABASE_PASSWORD")
			}
			
			dsn = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", 
				host, port, dbname, user, password)
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Erro ao conectar no banco: %v", err)
		}

		dbInstance = &postgresDatabase{Db: db}
	})

	return dbInstance
}

func (p *postgresDatabase) GetDb() *gorm.DB {
	return p.Db
}
