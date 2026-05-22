package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/utils"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	role := "staff"
	if id%2 == 0 {
		role = "visitor"
	}

	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       id,
		Name:     "User 1",
		Email:    fmt.Sprintf("%s%d@library.com", role, id),
		Password: hashedPassword,
		Role:     role,
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	id := rand.IntN(100)
	role := "staff"
	if len(email)%2 == 0 {
		role = "visitor"
	}

	hashedPassword, err := utils.HashPassword("password")
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       id,
		Name:     fmt.Sprintf("User %d", id),
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}

	return user, nil
}

func (r *UserRepository) GetUsersByRole(ctx context.Context, role string) ([]model.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM users WHERE role = ?`, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO users (id, name, email, password, role) VALUES (?, ?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.Password, user.Role)
	return err
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET name = ?, email = ?, role = ?, password = ? WHERE id = ?`,
		user.Name, user.Email, user.Role, user.Password, user.ID)
	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	return err
}
