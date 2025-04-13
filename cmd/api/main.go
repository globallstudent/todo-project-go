package main

import (
	"log"
	// "net/http"

	"github.com/gin-gonic/gin"
	"github.com/globallstudent/todo-project-go/internal/config"
	"github.com/globallstudent/todo-project-go/internal/database"
	"github.com/globallstudent/todo-project-go/internal/handlers"
	"github.com/globallstudent/todo-project-go/internal/middleware"
	"github.com/globallstudent/todo-project-go/internal/repositories"
	"github.com/globallstudent/todo-project-go/internal/services"
	// Swagger docs
	_ "github.com/globallstudent/todo-project-go/docs" // Import generated docs
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Todo API
// @version 1.0
// @description A simple Todo API with user authentication and role-based access control.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	todoRepo := repositories.NewTodoRepository(db)

	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	todoService := services.NewTodoService(todoRepo)

	authHandler := handlers.NewAuthHandler(authService)
	todoHandler := handlers.NewTodoHandler(todoService)

	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	protected := r.Group("/todos")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.POST("", todoHandler.CreateTodo)
		protected.GET("/:id", todoHandler.GetTodo)
		protected.GET("", todoHandler.GetTodos)
		protected.PUT("/:id", todoHandler.UpdateTodo)
		protected.DELETE("/:id", todoHandler.DeleteTodo)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	admin := r.Group("/admin/todos")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret), middleware.AdminMiddleware())
	{
		admin.GET("", todoHandler.GetTodos)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
