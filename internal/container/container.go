package container

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/config"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/controller"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/database"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type Container struct {
	AuthMenu    *AuthMenu
	VisitorMenu *VisitorMenu
	StaffMenu   *StaffMenu
}

func New() *Container {
	cfg := config.InitConfig()

	db, err := database.InitializeDB(cfg)
	if err != nil {
		fmt.Println("❌ Gagal terhubung ke database:", err)
		return nil
	}
	defer db.Close()

	// Repositories
	userRepository := repository.NewUserRepository(db)

	// Services
	authService := service.NewAuthService(userRepository)

	// Controllers
	authController := controller.NewAuthController(authService)

	return &Container{
		AuthMenu:    &AuthMenu{authController: authController},
		VisitorMenu: &VisitorMenu{},
		StaffMenu:   &StaffMenu{},
	}
}

func (c *Container) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("\n👋 Selamat tinggal!")
			return
		default:
		}

		user := c.AuthMenu.Main(ctx)
		if user == nil {
			fmt.Println("\n👋 Selamat tinggal!")
			return
		}

		switch user.Role {
		case "staff":
			c.StaffMenu.Dashboard(ctx)
		case "visitor":
			c.VisitorMenu.Dashboard(ctx)
		}
	}
}
