package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/controller"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/manifoldco/promptui"
)

type AuthMenu struct {
	authController      *controller.AuthController
	activityController  *controller.ActivityController
}

func (m *AuthMenu) Main(ctx context.Context) *model.User {
	prompt := promptui.Select{
		Label:        "LIBRARY CLI",
		Items:        []string{"Login", "Keluar Aplikasi"},
		Size:         10,
		HideSelected: true,
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
	for {
		prompt := promptui.Prompt{
			Label: "Email",
			Validate: func(input string) error {
				if strings.TrimSpace(input) == "" {
					return fmt.Errorf("Email tidak boleh kosong")
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
					return fmt.Errorf("Password tidak boleh kosong")
				}
				return nil
			},
		}
		password, err := prompt.Run()
		if err != nil {
			return nil
		}

		user, err := m.authController.Login(ctx, email, password)
		if err != nil {
			fmt.Printf("\n❌ %v\n\n", err)
			continue
		}

		m.activityController.Log(ctx, "Login", fmt.Sprintf("%s login sebagai %s", user.Name, user.Role))
		fmt.Printf("\n✅ Login berhasil — selamat datang, %s!\n\n", user.Name)
		return user
	}
}
