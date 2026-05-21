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
	authorRepository := repository.NewAuthorRepository(db)
	bookRepository := repository.NewBookRepository(db)

	// Services
	authService := service.NewAuthService(userRepository)
	authorService := service.NewAuthorService(authorRepository)
	bookService := service.NewBookService(bookRepository, authorRepository)

	// Controllers
	authController := controller.NewAuthController(authService)
	authorController := controller.NewAuthorController(authorService)
	bookController := controller.NewBookController(bookService)

	return &Container{
		AuthMenu:    &AuthMenu{authController: authController},
		VisitorMenu: &VisitorMenu{authorController: authorController, bookController: bookController},
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
