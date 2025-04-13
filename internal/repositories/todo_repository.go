package repositories

import (
	"context"

	"github.com/globallstudent/todo-project-go/internal/database"
	"github.com/globallstudent/todo-project-go/internal/models"
)

type TodoRepository struct {
	db *database.DB
}

func NewTodoRepository(db *database.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) CreateTodo(ctx context.Context, todo *models.Todo) error {
	query := `
		INSERT INTO todos (title, description, completed, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.db.Pool.QueryRow(ctx, query, todo.Title, todo.Description, todo.Completed, todo.UserID).
		Scan(&todo.ID, &todo.CreatedAt)
}

func (r *TodoRepository) FindTodoByID(ctx context.Context, id int) (*models.Todo, error) {
	todo := &models.Todo{}
	query := `
		SELECT id, title, description, completed, user_id, created_at
		FROM todos
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.UserID, &todo.CreatedAt)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *TodoRepository) FindTodosByUserID(ctx context.Context, userID int) ([]*models.Todo, error) {
	query := `
		SELECT id, title, description, completed, user_id, created_at
		FROM todos
		WHERE user_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		todo := &models.Todo{}
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.UserID, &todo.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepository) UpdateTodo(ctx context.Context, todo *models.Todo) error {
	query := `
		UPDATE todos
		SET title = $1, description = $2, completed = $3
		WHERE id = $4
	`
	_, err := r.db.Pool.Exec(ctx, query, todo.Title, todo.Description, todo.Completed, todo.ID)
	return err
}

func (r *TodoRepository) DeleteTodo(ctx context.Context, id int) error {
	query := `DELETE FROM todos WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
