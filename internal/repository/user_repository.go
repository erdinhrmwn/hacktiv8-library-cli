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

func (r *UserRepository) FindByID(ctx context.Context, id int) (*model.User, error) {
	// User Dummy
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

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	// User Dummy
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

func (r *UserRepository) Save(ctx context.Context, user *model.User) error {
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	return nil
}
