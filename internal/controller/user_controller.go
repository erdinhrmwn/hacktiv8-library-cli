package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) GetByRole(ctx context.Context, role string) ([]*model.User, error) {
	if role != "staff" && role != "visitor" {
		return nil, fmt.Errorf("role tidak valid")
	}
	return c.userService.GetUserByRole(ctx, role)
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
	Role     string
}

func (c *UserController) Create(ctx context.Context, input CreateUserInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("nama tidak boleh kosong")
	}
	if strings.TrimSpace(input.Email) == "" {
		return fmt.Errorf("email tidak boleh kosong")
	}
	if !strings.Contains(input.Email, "@") {
		return fmt.Errorf("email tidak valid")
	}
	if input.Password == "" {
		return fmt.Errorf("password tidak boleh kosong")
	}
	if len(input.Password) < 5 {
		return fmt.Errorf("password minimal 5 karakter")
	}
	if input.Role != "staff" && input.Role != "visitor" {
		return fmt.Errorf("role tidak valid")
	}

	user := &model.User{
		Name:     strings.TrimSpace(input.Name),
		Email:    strings.TrimSpace(input.Email),
		Password: input.Password,
		Role:     input.Role,
	}
	return c.userService.CreateUser(ctx, user)
}

type UpdateUserInput struct {
	UserID int
	Name   string
	Email  string
	Role   string
}

func (c *UserController) Update(ctx context.Context, input UpdateUserInput) error {
	if input.UserID <= 0 {
		return fmt.Errorf("ID user tidak valid")
	}
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("nama tidak boleh kosong")
	}
	if strings.TrimSpace(input.Email) == "" {
		return fmt.Errorf("email tidak boleh kosong")
	}
	if !strings.Contains(input.Email, "@") {
		return fmt.Errorf("email tidak valid")
	}
	if input.Role != "staff" && input.Role != "visitor" {
		return fmt.Errorf("role tidak valid")
	}

	user := &model.User{
		ID:    input.UserID,
		Name:  strings.TrimSpace(input.Name),
		Email: strings.TrimSpace(input.Email),
		Role:  input.Role,
	}
	return c.userService.UpdateUser(ctx, user)
}

type ChangePasswordInput struct {
	UserID      int
	OldPassword string
	NewPassword string
}

func (c *UserController) ChangePassword(ctx context.Context, input ChangePasswordInput) error {
	if input.UserID <= 0 {
		return fmt.Errorf("ID user tidak valid")
	}
	if strings.TrimSpace(input.OldPassword) == "" {
		return fmt.Errorf("password lama tidak boleh kosong")
	}
	if strings.TrimSpace(input.NewPassword) == "" {
		return fmt.Errorf("password baru tidak boleh kosong")
	}
	if len(input.NewPassword) < 5 {
		return fmt.Errorf("password baru minimal 5 karakter")
	}
	if input.OldPassword == input.NewPassword {
		return fmt.Errorf("password baru tidak boleh sama dengan password lama")
	}

	return c.userService.ChangePassword(ctx, input.UserID, input.OldPassword, input.NewPassword)
}

func (c *UserController) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("ID user tidak valid")
	}
	return c.userService.DeleteUser(ctx, id)
}
