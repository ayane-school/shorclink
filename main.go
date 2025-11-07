package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"shorclick/handlers"
	"shorclick/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, proceeding with environment variables")
	}
	original_link := os.Getenv("ORIGINAL_LINK")
	// DB接続の初期化
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := "5432"
	log.Println("Connecting to database")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	log.Println("Database connection established")

	// マイグレーションの実行
	err = db.AutoMigrate(&models.ShortLink{})
	if err != nil {
		log.Fatalln("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")

	// Ginのセットアップ
	r:= gin.Default()
	log.Println("Starting server on :8080")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, original_link)
	})
	r.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "don't setup api token",
		})
	})
	r.GET("/:id", handlers.ShortLink(db))
	r.Run(":8080")
}