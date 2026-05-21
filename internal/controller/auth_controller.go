package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Login(ctx context.Context, email, password string) (*model.User, error) {
	if strings.TrimSpace(email) == "" {
		return nil, fmt.Errorf("Email tidak boleh kosong")
	}
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("Email tidak valid")
	}
	if strings.TrimSpace(password) == "" {
		return nil, fmt.Errorf("Password tidak boleh kosong")
	}
	if len(password) < 5 {
		return nil, fmt.Errorf("Password minimal 5 karakter")
	}

	return c.authService.Login(ctx, email, password)
}
