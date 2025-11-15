package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(c *gin.Context) {
	var loginRequest LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}
	if loginRequest.Username != "admin" || loginRequest.Password != "password" {
		ResponseJSON(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		ResponseJSON(c, http.StatusInternalServerError, "Could not generate token", nil)
		return
	}
	ResponseJSON(c, http.StatusOK, "Token generated successfully", gin.H{"token": tokenString})
}

func CreateBook(c *gin.Context) {
	var book Book

	if err := c.ShouldBindJSON(&book); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	DB.Create(&book)
	ResponseJSON(c, http.StatusCreated, "Book created successfully", book)
}

func GetBooks(c *gin.Context) {
	var books []Book
	DB.Find(&books)
	ResponseJSON(c, http.StatusOK, "Books retrieved successfully", books)
}

func GetBook(c *gin.Context) {
	var book Book
	id := c.Param("id")

	if err := DB.First(&book, id).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}

	ResponseJSON(c, http.StatusOK, "Book retrieved successfully", book)
}

func UpdateBook(c *gin.Context) {
	var book Book
	id := c.Param("id")

	if err := DB.First(&book, id).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		ResponseJSON(c, http.StatusBadRequest, "Invalid input", nil)
		return
	}

	DB.Save(&book)
	ResponseJSON(c, http.StatusOK, "Book updated successfully", book)
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")

	if err := DB.Delete(&Book{}, id).Error; err != nil {
		ResponseJSON(c, http.StatusNotFound, "Book not found", nil)
		return
	}

	ResponseJSON(c, http.StatusOK, "Book deleted successfully", nil)
}
