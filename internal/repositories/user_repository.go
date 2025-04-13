package repositories

import (
	"context"

	"github.com/globallstudent/todo-project-go/internal/database"
	"github.com/globallstudent/todo-project-go/internal/models"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, password, role)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.Pool.QueryRow(ctx, query, user.Username, user.Password, user.Role).
		Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, password, role, created_at
		FROM users
		WHERE username = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
