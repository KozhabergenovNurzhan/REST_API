package handlers

import (
	"go_book_api/cmd/api"
	"go_book_api/internal/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(c *gin.Context) {
	var loginRequest main.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		main.ResponseJSON(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}
	if loginRequest.Username != "admin" || loginRequest.Password != "password" {
		main.ResponseJSON(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(main.jwtSecret)
	if err != nil {
		main.ResponseJSON(c, http.StatusInternalServerError, "Could not generate token", nil)
		return
	}
	main.ResponseJSON(c, http.StatusOK, "Token generated successfully", gin.H{"token": tokenString})
}

func CreateBook(c *gin.Context) {
	var book main.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		main.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	database.DB.Create(&book)
	main.ResponseJSON(c, http.StatusCreated, "Book created successfully", book)
}

func GetBooks(c *gin.Context) {
	var books []main.Book
	database.DB.Find(&books)
	main.ResponseJSON(c, http.StatusOK, "Books retrieved successfully", books)
}

func GetBook(c *gin.Context) {
	var book main.Book
	id := c.Param("id")

	if err := database.DB.First(&book, id).Error; err != nil {
		main.ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}

	main.ResponseJSON(c, http.StatusOK, "Book retrieved successfully", book)
}

func UpdateBook(c *gin.Context) {
	var book main.Book
	id := c.Param("id")

	if err := database.DB.First(&book, id).Error; err != nil {
		main.ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		main.ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	database.DB.Save(&book)
	main.ResponseJSON(c, http.StatusOK, "Book updated successfully", book)
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&main.Book{}, id).Error; err != nil {
		main.ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}

	main.ResponseJSON(c, http.StatusOK, "Book deleted successfully", nil)
}
