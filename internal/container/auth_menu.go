package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/manifoldco/promptui"
)

type AuthMenu struct{}

func (m *AuthMenu) Main(ctx context.Context) *model.User {
	prompt := promptui.Select{
		Label: "LIBRARY CLI",
		Items: []string{"Login", "Keluar Aplikasi"},
		Size:  10,
	}
	_, sel, err := prompt.Run()
	if err != nil {
		return nil
	}

	switch sel {
	case "Login":
		return m.login(ctx)
	default:
		return nil
	}
}

func (m *AuthMenu) login(ctx context.Context) *model.User {
	prompt := promptui.Prompt{
		Label: "Email",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("email tidak boleh kosong")
			}
			return nil
		},
	}
	email, err := prompt.Run()
	if err != nil {
		return nil
	}

	prompt = promptui.Prompt{
		Label: "Password",
		Mask:  '*',
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("password tidak boleh kosong")
			}
			return nil
		},
	}
	_, err = prompt.Run()
	if err != nil {
		return nil
	}

	fmt.Printf("\n✅ Login berhasil — selamat datang, %s!\n\n", email)

	// TODO: authenticate via AuthService
	return &model.User{
		Email: email,
		Role:  "staff",
	}
}
