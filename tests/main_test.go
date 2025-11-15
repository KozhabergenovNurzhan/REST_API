package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"go_book_api/api"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("SECRET_TOKEN"))

func setupTestDB() {
	var err error
	api.DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}
	api.DB.AutoMigrate(&api.Book{})
}

func addBook() api.Book {
	book := api.Book{Title: "Go programming", Author: "John Doe", Year: 2025}
	api.DB.Create(&book)
	return book
}

func generateValidToken() string {
	expirationTime := time.Now().Add(15 * time.Minute)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": expirationTime.Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

func TestGenerateJWT(t *testing.T) {
	router := gin.Default()
	router.POST("/token", api.GenerateJWT)
	
	loginRequest := map[string]string{
		"username": "admin",
		"password": "password",
	}

	jsonValue, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/token", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, status)
	}

	var response api.JsonResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Data == nil || response.Data.(map[string]interface{})["token"] == "" {
		t.Errorf("Expected token in response, got nil or empty")
	}
}

func TestCreateBook(t *testing.T) {
	setupTestDB()
	router := gin.Default()
	router.POST("/book", api.CreateBook)

	newBook := api.Book{Title: "Demo Book", Author: "Demo Author", Year: 2021}
	jsonValue, _ := json.Marshal(newBook)

	req, _ := http.NewRequest("POST", "/book", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestGetBooks(t *testing.T) {
	setupTestDB()
	addBook()

	router := gin.Default()
	router.GET("/books", api.GetBooks)

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetBook(t *testing.T) {
	setupTestDB()
	book := addBook()

	router := gin.Default()
	router.GET("/book/:id", api.GetBook)

	req, _ := http.NewRequest("GET", "/book/"+strconv.Itoa(int(book.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateBook(t *testing.T) {
	setupTestDB()
	book := addBook()

	router := gin.Default()
	router.PUT("/book/:id", api.UpdateBook)

	updateBook := api.Book{
		Title:  "Advanced Go Programming",
		Author: "Demo Author",
		Year:   2025,
	}

	body, _ := json.Marshal(updateBook)
	req, _ := http.NewRequest("PUT",
		"/book/"+strconv.Itoa(int(book.ID)),
		bytes.NewBuffer(body))

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteBook(t *testing.T) {
	setupTestDB()
	book := addBook()

	router := gin.Default()
	router.DELETE("/book/:id", api.DeleteBook)

	req, _ := http.NewRequest("DELETE",
		"/book/"+strconv.Itoa(int(book.ID)),
		nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}

	var deletedBook api.Book
	err := api.DB.First(&deletedBook, book.ID).Error
	if err == nil {
		t.Errorf("Expected book to be deleted, but still exists")
	}
}
