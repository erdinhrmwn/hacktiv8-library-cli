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
	activityRepository := repository.NewActivityRepository(db)
	loanRepository := repository.NewLoanRepository(db)
	invoiceRepository := repository.NewInvoiceRepository(db)
	paymentRepository := repository.NewPaymentRepository(db)

	// Services
	authService := service.NewAuthService(userRepository)
	authorService := service.NewAuthorService(authorRepository)
	bookService := service.NewBookService(bookRepository, authorRepository)
	userService := service.NewUserService(userRepository)
	activityService := service.NewActivityService(activityRepository)
	loanService := service.NewLoanService(loanRepository, bookRepository, invoiceRepository)
	invoiceService := service.NewInvoiceService(invoiceRepository)
	paymentService := service.NewPaymentService(paymentRepository, invoiceRepository)

	// Controllers
	authController := controller.NewAuthController(authService)
	authorController := controller.NewAuthorController(authorService)
	bookController := controller.NewBookController(bookService)
	userController := controller.NewUserController(userService)
	activityController := controller.NewActivityController(activityService)
	loanController := controller.NewLoanController(loanService)
	invoiceController := controller.NewInvoiceController(invoiceService)
	paymentController := controller.NewPaymentController(paymentService)

	return &Container{
		AuthMenu: &AuthMenu{
			authController:     authController,
			activityController: activityController,
		},
		VisitorMenu: &VisitorMenu{
			authorController:   authorController,
			bookController:     bookController,
			userController:     userController,
			activityController: activityController,
			loanController:     loanController,
			invoiceController:  invoiceController,
		},
		StaffMenu: &StaffMenu{
			authorController:   authorController,
			bookController:     bookController,
			userController:     userController,
			activityController: activityController,
			loanController:     loanController,
			invoiceController:  invoiceController,
			paymentController:  paymentController,
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
