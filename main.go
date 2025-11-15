package main

import (
	"log"

	"go_book_api/api"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	api.InitDB()

	router := gin.Default()
	router.POST("/token", api.GenerateJWT)

	if err := api.DB.AutoMigrate(&api.Book{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	protected := router.Group("/", api.JWTAuthMiddleware())
	{
		protected.POST("/book", api.CreateBook)
		protected.GET("/books", api.GetBooks)
		protected.GET("/book/:id", api.GetBook)
		protected.PUT("/book/:id", api.UpdateBook)
		protected.DELETE("/book/:id", api.DeleteBook)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
