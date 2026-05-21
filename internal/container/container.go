package container

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/controller"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type Container struct {
	AuthMenu    *AuthMenu
	VisitorMenu *VisitorMenu
	StaffMenu   *StaffMenu
}

func New(db *sql.DB) *Container {
	// Repositories
	userRepository := repository.NewUserRepository(db)
	authorRepository := repository.NewAuthorRepository(db)
	bookRepository := repository.NewBookRepository(db)

	// Services
	authService := service.NewAuthService(userRepository)
	authorService := service.NewAuthorService(authorRepository)
	bookService := service.NewBookService(bookRepository, authorRepository)
	userService := service.NewUserService(userRepository)

	// Controllers
	authController := controller.NewAuthController(authService)
	authorController := controller.NewAuthorController(authorService)
	bookController := controller.NewBookController(bookService)
	userController := controller.NewUserController(userService)

	return &Container{
		AuthMenu: &AuthMenu{
			authController: authController,
		},
		VisitorMenu: &VisitorMenu{
			authorController: authorController,
			bookController:   bookController,
			userController:   userController,
		},
		StaffMenu: &StaffMenu{
			authorController: authorController,
			bookController:   bookController,
			userController:   userController,
		},
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
			c.StaffMenu.Dashboard(ctx, user)
		case "visitor":
			c.VisitorMenu.Dashboard(ctx, user)
		}
	}
}
