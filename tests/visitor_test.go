package tests

import (
	"context"
	"testing"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/controller"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
	"github.com/erdinhrmwn/hacktiv8-library-cli/utils"
)

func TestVisitorFlow(t *testing.T) {
	ctx := context.Background()

	// repositories
	userRepo := repository.NewUserRepository(testDB)
	bookRepo := repository.NewBookRepository(testDB)
	authorRepo := repository.NewAuthorRepository(testDB)
	loanRepo := repository.NewLoanRepository(testDB)
	invoiceRepo := repository.NewInvoiceRepository(testDB)

	// services
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	bookService := service.NewBookService(bookRepo, authorRepo)
	loanService := service.NewLoanService(loanRepo, bookRepo, invoiceRepo)
	invoiceService := service.NewInvoiceService(invoiceRepo)

	// controllers
	authCtrl := controller.NewAuthController(authService)
	userCtrl := controller.NewUserController(userService)
	bookCtrl := controller.NewBookController(bookService)
	loanCtrl := controller.NewLoanController(loanService)
	invoiceCtrl := controller.NewInvoiceController(invoiceService)

	var visitorID int

	// Login
	t.Run("Login", func(t *testing.T) {
		user, err := authCtrl.Login(ctx, "visitor@test.com", "password")
		if err != nil {
			t.Fatal("gagal login:", err)
		}
		visitorID = user.ID
		if user.Role != "visitor" {
			t.Fatal("role should be visitor")
		}
	})

	// View all books
	t.Run("View all books", func(t *testing.T) {
		books, err := bookCtrl.GetAllBooks(ctx)
		if err != nil {
			t.Fatal("gagal lihat buku:", err)
		}
		if len(books) < 3 {
			t.Fatal("seharusnya >= 3 buku, got:", len(books))
		}
	})

	// Search books
	t.Run("Search books", func(t *testing.T) {
		books, err := bookCtrl.SearchBook(ctx, "Book")
		if err != nil {
			t.Fatal("gagal cari buku:", err)
		}
		if len(books) == 0 {
			t.Fatal("seharusnya ada hasil pencarian")
		}
	})

	// View book detail
	t.Run("View book detail", func(t *testing.T) {
		books, _ := bookCtrl.GetAllBooks(ctx)
		book, err := bookCtrl.GetBookByID(ctx, books[0].ID)
		if err != nil {
			t.Fatal("gagal lihat detail:", err)
		}
		if book.Title == "" {
			t.Fatal("judul buku kosong")
		}
	})

	// My Loans
	t.Run("Check loans", func(t *testing.T) {
		loans, err := loanCtrl.GetActiveByVisitorID(ctx, visitorID)
		if err != nil {
			t.Fatal("gagal ambil loans:", err)
		}
		if len(loans) != 0 {
			t.Fatal("seharusnya 0 loans, got:", len(loans))
		}
	})

	// My Invoices
	t.Run("Check invoices", func(t *testing.T) {
		invoices, err := invoiceCtrl.GetUnpaidByUserID(ctx, visitorID)
		if err != nil {
			t.Fatal("gagal ambil invoices:", err)
		}
		if len(invoices) != 0 {
			t.Fatal("seharusnya 0 invoices, got:", len(invoices))
		}
	})

	// Change name
	t.Run("Change name", func(t *testing.T) {
		err := userCtrl.Update(ctx, controller.UpdateUserInput{
			UserID: visitorID,
			Name:   "Visitor Updated",
			Email:  "visitor-updated@test.com",
			Role:   "visitor",
		})
		if err != nil {
			t.Fatal("gagal ubah nama:", err)
		}

		user, _ := userRepo.GetUserByID(ctx, visitorID)
		if user.Name != "Visitor Updated" {
			t.Fatal("nama belum berubah, got:", user.Name)
		}
	})

	// Change password
	t.Run("Change password", func(t *testing.T) {
		err := userCtrl.ChangePassword(ctx, controller.ChangePasswordInput{
			UserID:      visitorID,
			OldPassword: "password",
			NewPassword: "newpass123",
		})
		if err != nil {
			t.Fatal("gagal ganti password:", err)
		}

		user, _ := userRepo.GetUserByID(ctx, visitorID)
		if !utils.VerifyPassword("newpass123", user.Password) {
			t.Fatal("password baru tidak cocok")
		}
	})
}
