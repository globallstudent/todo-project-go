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
)

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
