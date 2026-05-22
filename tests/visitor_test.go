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

	userRepo := repository.NewUserRepository(testDB)
	authService := service.NewAuthService(userRepo)
	authCtrl := controller.NewAuthController(authService)
	userService := service.NewUserService(userRepo)
	userCtrl := controller.NewUserController(userService)
	loanRepo := repository.NewLoanRepository(testDB)
	loanCtrl := controller.NewLoanController(nil, loanRepo)
	invoiceRepo := repository.NewInvoiceRepository(testDB)
	invoiceCtrl := controller.NewInvoiceController(invoiceRepo, loanRepo)

	t.Run("Login as visitor", func(t *testing.T) {
		user, err := authCtrl.Login(ctx, "visitor@test.com", "password")
		if err != nil {
			t.Fatal("gagal login:", err)
		}
		if user.Role != "visitor" {
			t.Fatal("role should be visitor, got:", user.Role)
		}
	})

	t.Run("Change name", func(t *testing.T) {
		err := userCtrl.Update(ctx, controller.UpdateUserInput{
			ID:    2,
			Name:  "Visitor Updated",
			Email: "visitor@test.com",
			Role:  "visitor",
		})
		if err != nil {
			t.Fatal("gagal ubah nama:", err)
		}

		user, _ := userRepo.GetUserByID(ctx, 2)
		if user.Name != "Visitor Updated" {
			t.Fatal("nama belum berubah, got:", user.Name)
		}
	})

	t.Run("Change password and verify", func(t *testing.T) {
		err := userCtrl.ChangePassword(ctx, controller.ChangePasswordInput{
			UserID:      2,
			OldPassword: "password",
			NewPassword: "newpass123",
		})
		if err != nil {
			t.Fatal("gagal ganti password:", err)
		}

		user, _ := userRepo.GetUserByID(ctx, 2)
		if !utils.VerifyPassword("newpass123", user.Password) {
			t.Fatal("password baru tidak cocok")
		}
	})

	t.Run("Check loans", func(t *testing.T) {
		loans, err := loanCtrl.GetActiveByVisitorID(ctx, 2)
		if err != nil {
			t.Fatal("gagal ambil loans:", err)
		}
		if len(loans) != 0 {
			t.Fatal("seharusnya 0 loans, got:", len(loans))
		}
	})

	t.Run("Check invoices", func(t *testing.T) {
		invoices, err := invoiceCtrl.GetUnpaidByUserID(ctx, 2)
		if err != nil {
			t.Fatal("gagal ambil invoices:", err)
		}
		if len(invoices) != 0 {
			t.Fatal("seharusnya 0 invoices, got:", len(invoices))
		}
	})
}
