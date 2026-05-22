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
	authController     *controller.AuthController
	userController     *controller.UserController
	activityController *controller.ActivityController
}

func (m *AuthMenu) Main(ctx context.Context) *model.User {
	prompt := promptui.Select{
		Label:        "LIBRARY CLI",
		Items:        []string{"Login", "Daftar", "Keluar Aplikasi"},
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
	case "Daftar":
		m.register(ctx)
		return nil
	default:
		return nil
	}
}

func (m *AuthMenu) register(ctx context.Context) {
	name, err := promptText("Nama", notEmpty)
	if err != nil {
		return
	}
	email, err := promptText("Email", validEmail)
	if err != nil {
		return
	}
	password, err := promptMasked("Password", minLength(5))
	if err != nil {
		return
	}

	err = m.userController.Create(ctx, controller.CreateUserInput{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     "visitor",
	})
	if err != nil {
		fmt.Printf("\n❌ Gagal mendaftar: %v\n\n", err)
		return
	}

	go m.activityController.Log(context.Background(), "Register", fmt.Sprintf("%s mendaftar sebagai visitor", email))
	fmt.Printf("\n✅ Pendaftaran berhasil! Silakan login.\n\n")
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

		go m.activityController.Log(context.Background(), "Login", fmt.Sprintf("%s login sebagai %s", user.Name, user.Role))
		fmt.Printf("\n✅ Login berhasil — selamat datang, %s!\n\n", user.Name)
		return user
	}
}

func promptText(label string, validate promptui.ValidateFunc) (string, error) {
	return (&promptui.Prompt{Label: label, Validate: validate}).Run()
}

func promptMasked(label string, validate promptui.ValidateFunc) (string, error) {
	return (&promptui.Prompt{Label: label, Mask: '*', Validate: validate}).Run()
}

func notEmpty(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("tidak boleh kosong")
	}
	return nil
}

func validEmail(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("email tidak boleh kosong")
	}
	if !strings.Contains(input, "@") {
		return fmt.Errorf("email tidak valid")
	}
	return nil
}

func minLength(n int) promptui.ValidateFunc {
	return func(input string) error {
		if strings.TrimSpace(input) == "" {
			return fmt.Errorf("tidak boleh kosong")
		}
		if len(input) < n {
			return fmt.Errorf("minimal %d karakter", n)
		}
		return nil
	}
}
