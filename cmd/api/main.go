package main

import (
	"knovel/internal/auth"
	"knovel/internal/database"
	"knovel/internal/handlers"
	"knovel/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := "host=postgres user=postgres password=postgres dbname=knovel port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	userRepo := repository.NewGormUserRepository(db)
	taskRepo := repository.NewGormTaskRepository(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, "knovel-secret-key")
	taskHandler := handlers.NewTaskHandler(taskRepo)

	// Initialize auth middleware
	authMiddleware := auth.NewAuthMiddleware("knovel-secret-key")

	// Setup router
	r := gin.Default()

	// Public routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(authMiddleware.AuthRequired())
	{
		tasks := api.Group("/tasks")
		{
			tasks.POST("/", taskHandler.CreateTask)
			tasks.GET("/", taskHandler.GetTasks)
		}
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
