package services

import (
	"context"
	"errors"

	"github.com/globallstudent/todo-project-go/internal/models"
	"github.com/globallstudent/todo-project-go/internal/repositories"
)

type TodoService struct {
	todoRepo *repositories.TodoRepository
}

func NewTodoService(todoRepo *repositories.TodoRepository) *TodoService {
	return &TodoService{todoRepo: todoRepo}
}

func (s *TodoService) CreateTodo(ctx context.Context, todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("title is required")
	}
	return s.todoRepo.CreateTodo(ctx, todo)
}

func (s *TodoService) GetTodoByID(ctx context.Context, id int) (*models.Todo, error) {
	return s.todoRepo.FindTodoByID(ctx, id)
}

func (s *TodoService) GetTodosByUserID(ctx context.Context, userID int) ([]*models.Todo, error) {
	return s.todoRepo.FindTodosByUserID(ctx, userID)
}

func (s *TodoService) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("title is required")
	}
	return s.todoRepo.UpdateTodo(ctx, todo)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
	return s.todoRepo.DeleteTodo(ctx, id)
}
