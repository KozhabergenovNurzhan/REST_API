package main

import (
	"go_book_api/internal/database"
	"go_book_api/internal/handlers"
	"go_book_api/internal/middleware"
	main3 "go_book_api/internal/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	database.InitDB()

	router := gin.Default()
	router.POST("/token", handlers.GenerateJWT)

	if err := database.DB.AutoMigrate(&main3.Book{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	protected := router.Group("/", middleware.JWTAuthMiddleware())
	{
		protected.POST("/book", handlers.CreateBook)
		protected.GET("/books", handlers.GetBooks)
		protected.GET("/book/:id", handlers.GetBook)
		protected.PUT("/book/:id", handlers.UpdateBook)
		protected.DELETE("/book/:id", handlers.DeleteBook)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
