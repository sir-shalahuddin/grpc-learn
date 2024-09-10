package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sir-shalahuddin/grpc-learn/userservice/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (name, email, password_hash) 
              VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.Password)
	if err != nil {
		log.Printf("[Repository - CreateUser] Error executing query: %v", err)
		return fmt.Errorf("failed to create user")
	}
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `SELECT id, name, email, created_at, updated_at, role FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userID)

	var user models.User
	err := row.Scan(&user.UserID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		log.Printf("[Repository - GetUserByID] Error scanning row: %v", err)
		return nil, fmt.Errorf("failed to retrieve user")
	}
	return &user, nil
}

// UpdateUser updates the user details in the database
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() 
              WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.UserID)
	if err != nil {
		log.Printf("[Repository - UpdateUser] Error executing query: %v", err)
		return fmt.Errorf("failed to update user")
	}
	return nil
}

// DeleteUser removes a user from the database
func (r *userRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		log.Printf("[Repository - DeleteUser] Error executing query: %v", err)
		return fmt.Errorf("failed to delete user")
	}
	return nil
}

func (r *userRepository) ListUsers(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, name, email, role, created_at, updated_at FROM users where role <> 'super admin'`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("[Repository - ListUsers] Error executing query: %v", err)
		return nil, fmt.Errorf("failed to list users")
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Printf("[Repository - ListUsers] Error scanning row: %v", err)
			return nil, fmt.Errorf("failed to list users")
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, name, email, password_hash FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)
	var user models.User
	err := row.Scan(&user.UserID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("[Repository - GetUserByEmail] Error scanning row: %v", err)
		return nil, fmt.Errorf("failed to retrieve user by email")
	}
	return &user, nil
}

func (r *userRepository) UpdateUserRoles(ctx context.Context, userID uuid.UUID, roles string) error {
	query := `UPDATE users SET role = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, roles, userID)
	if err != nil {
		log.Printf("[Repository - UpdateUserRoles] Error executing query: %v", err)
		return fmt.Errorf("failed to update user roles")
	}
	return nil
}
