package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"user-api/internal/database"
	"user-api/internal/handlers"
	"user-api/internal/middleware"
	"user-api/internal/repository"
	"user-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

//go:embed static/index.html
var indexHTML string

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbConfig := database.GetConfigFromEnv()
	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	router := gin.Default()

	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORS())

	// ГЛАВНАЯ СТРАНИЦА ИЗ ФАЙЛА static/index.html
	router.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(indexHTML))
	})

	// API Routes
	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("🌐 Web interface: http://localhost:%s", port)
	log.Printf("📡 API endpoint: http://localhost:%s/api/v1/users", port)
	log.Printf("❤️  Health check: http://localhost:%s/health", port)

	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
