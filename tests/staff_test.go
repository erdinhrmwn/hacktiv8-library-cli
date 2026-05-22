package tests

import (
	"context"
	"testing"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/controller"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

func TestStaffFlow(t *testing.T) {
	ctx := context.Background()

	userRepo := repository.NewUserRepository(testDB)
	authService := service.NewAuthService(userRepo)
	authCtrl := controller.NewAuthController(authService)
	authorRepo := repository.NewAuthorRepository(testDB)
	authorCtrl := controller.NewAuthorController(service.NewAuthorService(authorRepo))
	bookRepo := repository.NewBookRepository(testDB)
	bookCtrl := controller.NewBookController(service.NewBookService(bookRepo, authorRepo))
	loanRepo := repository.NewLoanRepository(testDB)
	loanService := service.NewLoanService(loanRepo, bookRepo, repository.NewInvoiceRepository(testDB))
	loanCtrl := controller.NewLoanController(loanService, loanRepo)
	invoiceRepo := repository.NewInvoiceRepository(testDB)
	invoiceCtrl := controller.NewInvoiceController(invoiceRepo, loanRepo)
	paymentRepo := repository.NewPaymentRepository(testDB)
	paymentCtrl := controller.NewPaymentController(paymentRepo, invoiceRepo)

	// Login as staff
	t.Run("Login as staff", func(t *testing.T) {
		user, err := authCtrl.Login(ctx, "staff@test.com", "password")
		if err != nil {
			t.Fatal("gagal login:", err)
		}
		if user.Role != "staff" {
			t.Fatal("role should be staff")
		}
	})

	// Lihat daftar author
	t.Run("List authors", func(t *testing.T) {
		authors, err := authorCtrl.GetAllAuthors(ctx)
		if err != nil {
			t.Fatal("gagal lihat author:", err)
		}
		if len(authors) < 2 {
			t.Fatal("seharusnya >= 2 author, got:", len(authors))
		}
	})

	// Create author
	var newAuthorID int
	t.Run("Create author", func(t *testing.T) {
		err := authorCtrl.AddAuthor(ctx, controller.AddAuthorInput{
			Name:        "Author Beta",
			Nationality: "Japanese",
			Bio:         "Bio beta",
		})
		if err != nil {
			t.Fatal("gagal buat author:", err)
		}

		authors, _ := authorCtrl.GetAllAuthors(ctx)
		for _, a := range authors {
			if a.Name == "Author Beta" {
				newAuthorID = a.ID
				break
			}
		}
		if newAuthorID == 0 {
			t.Fatal("author baru tidak ditemukan")
		}
	})

	// Lihat daftar book
	t.Run("List books", func(t *testing.T) {
		books, err := bookCtrl.GetAllBooks(ctx)
		if err != nil {
			t.Fatal("gagal lihat buku:", err)
		}
		if len(books) < 3 {
			t.Fatal("seharusnya >= 3 buku, got:", len(books))
		}
	})

	// Create book
	var newBookID int
	t.Run("Create book", func(t *testing.T) {
		err := bookCtrl.AddBook(ctx, controller.AddBookInput{
			ISBN:     "444-444",
			Title:    "Book Beta",
			Genre:    "Comedy",
			AuthorID: newAuthorID,
			Stock:    2,
		})
		if err != nil {
			t.Fatal("gagal buat buku:", err)
		}

		books, _ := bookCtrl.GetAllBooks(ctx)
		for _, b := range books {
			if b.Title == "Book Beta" {
				newBookID = b.ID
				break
			}
		}
		if newBookID == 0 {
			t.Fatal("buku baru tidak ditemukan")
		}
	})

	// Borrow book
	var loanID int
	t.Run("Borrow book", func(t *testing.T) {
		err := loanCtrl.BorrowBook(ctx, 2, 1, newBookID) // visitor=2, staff=1
		if err != nil {
			t.Fatal("gagal pinjam:", err)
		}

		loans, _ := loanCtrl.GetActiveByVisitorID(ctx, 2)
		if len(loans) != 1 {
			t.Fatal("seharusnya 1 loan aktif, got:", len(loans))
		}
		loanID = loans[0].ID

		// Cek stok berkurang
		book, _ := bookCtrl.GetBookByID(ctx, newBookID)
		if book.Stock != 1 {
			t.Fatal("stok harus 1, got:", book.Stock)
		}
	})

	// Return book before due - no invoice
	t.Run("Return book before due (no invoice)", func(t *testing.T) {
		fine, err := loanCtrl.ReturnBook(ctx, loanID)
		if err != nil {
			t.Fatal("gagal return:", err)
		}
		if fine != nil {
			t.Fatal("tidak seharusnya ada denda")
		}

		// Cek stok kembali
		book, _ := bookCtrl.GetBookByID(ctx, newBookID)
		if book.Stock != 2 {
			t.Fatal("stok harus 2, got:", book.Stock)
		}
	})

	// Borrow again, then return after due date
	t.Run("Borrow again for late return test", func(t *testing.T) {
		loanCtrl.BorrowBook(ctx, 2, 1, newBookID)
	})

	// Return after due date - invoice must be created
	t.Run("Return book after due date (invoice created)", func(t *testing.T) {
		loans, _ := loanCtrl.GetActiveByVisitorID(ctx, 2)
		if len(loans) == 0 {
			t.Fatal("tidak ada loan aktif")
		}

		// Set due date ke kemarin manual via DB
		testDB.Exec("UPDATE loans SET due_date = DATE_SUB(NOW(), INTERVAL 3 DAY) WHERE id = ?", loans[0].ID)

		fine, err := loanCtrl.ReturnBook(ctx, loans[0].ID)
		if err == nil {
			t.Fatal("seharusnya error karena telat")
		}
		if fine == nil {
			t.Fatal("seharusnya ada denda")
		}

		// Cek invoice
		invoices, _ := invoiceCtrl.GetUnpaidByUserID(ctx, 2)
		if len(invoices) == 0 {
			t.Fatal("seharusnya ada invoice")
		}
	})

	// Pay invoice
	t.Run("Pay invoice", func(t *testing.T) {
		invoices, _ := invoiceCtrl.GetUnpaidByUserID(ctx, 2)
		if len(invoices) == 0 {
			t.Fatal("tidak ada invoice untuk dibayar")
		}

		err := paymentCtrl.PayInvoice(ctx, invoices[0].ID, invoices[0].Amount, "cash")
		if err != nil {
			t.Fatal("gagal bayar:", err)
		}

		// Cek invoice sudah paid
		remaining, _ := invoiceCtrl.GetUnpaidByUserID(ctx, 2)
		if len(remaining) != 0 {
			t.Fatal("invoice seharusnya sudah paid")
		}
	})
}
