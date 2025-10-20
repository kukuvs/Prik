package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-api/internal/handlers"
	"user-api/internal/models"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestCreateUser(t *testing.T) {
	router := setupTestRouter()

	// Mock service
	mockService := &mockUserService{}
	handler := handlers.NewUserHandler(mockService)

	router.POST("/users", handler.CreateUser)

	user := models.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
	}

	jsonData, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetUser(t *testing.T) {
	router := setupTestRouter()

	mockService := &mockUserService{}
	handler := handlers.NewUserHandler(mockService)

	router.GET("/users/:id", handler.GetUser)

	req, _ := http.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUsers(t *testing.T) {
	router := setupTestRouter()

	mockService := &mockUserService{}
	handler := handlers.NewUserHandler(mockService)

	router.GET("/users", handler.GetUsers)

	req, _ := http.NewRequest("GET", "/users?page=1&page_size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateUser(t *testing.T) {
	router := setupTestRouter()

	mockService := &mockUserService{}
	handler := handlers.NewUserHandler(mockService)

	router.PUT("/users/:id", handler.UpdateUser)

	update := models.UpdateUserRequest{
		Name: "Updated Name",
	}

	jsonData, _ := json.Marshal(update)
	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUser(t *testing.T) {
	router := setupTestRouter()

	mockService := &mockUserService{}
	handler := handlers.NewUserHandler(mockService)

	router.DELETE("/users/:id", handler.DeleteUser)

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// Mock service
type mockUserService struct{}

func (m *mockUserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	return &models.User{
		ID:    1,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}, nil
}

func (m *mockUserService) GetUser(id int) (*models.User, error) {
	return &models.User{
		ID:    id,
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
	}, nil
}

func (m *mockUserService) GetUsers(page, pageSize int, filters map[string]interface{}) (*models.UserListResponse, error) {
	return &models.UserListResponse{
		Users: []models.User{
			{ID: 1, Name: "User 1", Email: "user1@example.com", Age: 25},
		},
		Total:      1,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: 1,
	}, nil
}

func (m *mockUserService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
	return &models.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}, nil
}

func (m *mockUserService) DeleteUser(id int) error {
	return nil
}
