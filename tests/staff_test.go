package tests

import (
	"context"
	"testing"
	"time"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/controller"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

func TestStaffFlow(t *testing.T) {
	ctx := context.Background()

	// repositories
	userRepo := repository.NewUserRepository(testDB)
	authorRepo := repository.NewAuthorRepository(testDB)
	bookRepo := repository.NewBookRepository(testDB)
	loanRepo := repository.NewLoanRepository(testDB)
	invoiceRepo := repository.NewInvoiceRepository(testDB)
	paymentRepo := repository.NewPaymentRepository(testDB)
	activityRepo := repository.NewActivityRepository(testDB)

	// services
	authService := service.NewAuthService(userRepo)
	authorService := service.NewAuthorService(authorRepo)
	bookService := service.NewBookService(bookRepo, authorRepo)
	loanService := service.NewLoanService(loanRepo, bookRepo, invoiceRepo)
	invoiceService := service.NewInvoiceService(invoiceRepo)
	paymentService := service.NewPaymentService(paymentRepo, invoiceRepo)
	userService := service.NewUserService(userRepo)
	activityService := service.NewActivityService(activityRepo)

	// controllers
	authCtrl := controller.NewAuthController(authService)
	authorCtrl := controller.NewAuthorController(authorService)
	bookCtrl := controller.NewBookController(bookService)
	loanCtrl := controller.NewLoanController(loanService)
	invoiceCtrl := controller.NewInvoiceController(invoiceService)
	paymentCtrl := controller.NewPaymentController(paymentService)
	userCtrl := controller.NewUserController(userService)
	activityCtrl := controller.NewActivityController(activityService)

	var staffID, visitorID int

	// Login
	t.Run("Login", func(t *testing.T) {
		user, err := authCtrl.Login(ctx, "staff@test.com", "password")
		if err != nil {
			t.Fatal("gagal login:", err)
		}
		staffID = user.ID
		if user.Role != "staff" {
			t.Fatal("role should be staff")
		}

		visitor, _ := userRepo.GetUserByEmail(ctx, "visitor@test.com")
		visitorID = visitor.ID
	})

	// List authors
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
			BirthDate:   time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC),
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

	// List books
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
		err := loanCtrl.BorrowBook(ctx, visitorID, staffID, newBookID)
		if err != nil {
			t.Fatal("gagal pinjam:", err)
		}

		loans, _ := loanCtrl.GetActiveByVisitorID(ctx, visitorID)
		if len(loans) != 1 {
			t.Fatal("seharusnya 1 loan aktif, got:", len(loans))
		}
		loanID = loans[0].ID

		// check stock
		book, _ := bookCtrl.GetBookByID(ctx, newBookID)
		if book.Stock != 1 {
			t.Fatal("stok harus 1, got:", book.Stock)
		}
	})

	// Return before due
	t.Run("Return book before due date (no invoice)", func(t *testing.T) {
		fine, err := loanCtrl.ReturnBook(ctx, loanID)
		if err != nil {
			t.Fatal("gagal return:", err)
		}
		if fine != nil {
			t.Fatal("tidak seharusnya ada denda")
		}

		// check stock
		book, _ := bookCtrl.GetBookByID(ctx, newBookID)
		if book.Stock != 2 {
			t.Fatal("stok harus 2, got:", book.Stock)
		}
	})

	// Return after due
	t.Run("Return book after due date (invoice created)", func(t *testing.T) {
		// borrow
		err := loanCtrl.BorrowBook(ctx, visitorID, staffID, newBookID)
		if err != nil {
			t.Fatal("gagal pinjam:", err)
		}

		loans, _ := loanCtrl.GetActiveByVisitorID(ctx, visitorID)
		if len(loans) == 0 {
			t.Fatal("tidak ada loan aktif")
		}

		// force late
		testDB.Exec("UPDATE loans SET due_date = DATE_SUB(NOW(), INTERVAL 3 DAY) WHERE id = ?", loans[0].ID)

		fine, err := loanCtrl.ReturnBook(ctx, loans[0].ID)
		if err == nil {
			t.Fatal("seharusnya error karena telat")
		}
		if fine == nil {
			t.Fatal("seharusnya ada denda")
		}

		// check invoice
		invoices, _ := invoiceCtrl.GetUnpaidByUserID(ctx, visitorID)
		if len(invoices) == 0 {
			t.Fatal("seharusnya ada invoice")
		}
	})

	// Pay invoice
	t.Run("Pay invoice", func(t *testing.T) {
		invoices, _ := invoiceCtrl.GetUnpaidByUserID(ctx, visitorID)
		if len(invoices) == 0 {
			t.Fatal("tidak ada invoice untuk dibayar")
		}

		err := paymentCtrl.PayInvoice(ctx, invoices[0].ID, invoices[0].Amount, "cash")
		if err != nil {
			t.Fatal("gagal bayar:", err)
		}

		// check invoice sudah paid
		remaining, _ := invoiceCtrl.GetUnpaidByUserID(ctx, visitorID)
		if len(remaining) != 0 {
			t.Fatal("invoice seharusnya sudah paid")
		}
	})

	// Register visitor
	t.Run("Register visitor", func(t *testing.T) {
		err := userCtrl.Create(ctx, controller.CreateUserInput{
			Name:     "New Visitor",
			Email:    "new@test.com",
			Password: "password",
			Role:     "visitor",
		})
		if err != nil {
			t.Fatal("gagal daftar visitor:", err)
		}

		users, _ := userCtrl.GetByRole(ctx, "visitor")
		found := false
		for _, u := range users {
			if u.Email == "new@test.com" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("visitor baru tidak ditemukan")
		}
	})

	// List visitors
	t.Run("List visitors", func(t *testing.T) {
		users, err := userCtrl.GetByRole(ctx, "visitor")
		if err != nil {
			t.Fatal("gagal lihat visitor:", err)
		}
		if len(users) < 2 {
			t.Fatal("seharusnya >= 2 visitor, got:", len(users))
		}
	})

	// Update book stock
	t.Run("Update stock", func(t *testing.T) {
		books, _ := bookCtrl.GetAllBooks(ctx)
		err := bookCtrl.UpdateBook(ctx, controller.UpdateBookInput{
			ID:       books[0].ID,
			ISBN:     books[0].ISBN,
			Title:    books[0].Title,
			Genre:    books[0].Genre,
			AuthorID: books[0].AuthorID,
			Stock:    10,
		})
		if err != nil {
			t.Fatal("gagal update stok:", err)
		}

		book, _ := bookCtrl.GetBookByID(ctx, books[0].ID)
		if book.Stock != 10 {
			t.Fatal("stok harus 10, got:", book.Stock)
		}
	})

	// Delete book
	t.Run("Delete book", func(t *testing.T) {
		books, _ := bookCtrl.GetAllBooks(ctx)
		// delete book with no loans
		for _, b := range books {
			if b.Title == "Book Two" {
				err := bookCtrl.DeleteBook(ctx, b.ID)
				if err != nil {
					t.Fatal("gagal hapus buku:", err)
				}
				return
			}
		}
	})

	// Activity logs
	t.Run("Activity logs", func(t *testing.T) {
		_, err := activityCtrl.GetAll(ctx)
		if err != nil {
			t.Fatal("gagal ambil log:", err)
		}
	})
}
