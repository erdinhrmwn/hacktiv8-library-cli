package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

// DendaPerHari adalah denda keterlambatan dalam satuan rupiah per hari.
const DendaPerHari = 5000.0

type LoanService struct {
	loanRepo     *repository.LoanRepository
	bookRepo     *repository.BookRepository
	invoiceRepo  *repository.InvoiceRepository
	activityRepo *repository.ActivityRepository
}

func NewLoanService(
	lr *repository.LoanRepository,
	br *repository.BookRepository,
	ir *repository.InvoiceRepository,
	ar *repository.ActivityRepository,
) *LoanService {
	return &LoanService{
		loanRepo:     lr,
		bookRepo:     br,
		invoiceRepo:  ir,
		activityRepo: ar,
	}
}

// Borrow memproses peminjaman buku oleh visitor yang dilayani staff.
// F-LOAN-01: kaitkan visitor_id, staff_id, book_id
// F-LOAN-03: tolak jika stok = 0
// F-LOAN-04: kurangi stok setelah berhasil
func (s *LoanService) Borrow(visitorID, staffID, bookID int) error {
	// 1. Ambil data buku & validasi stok
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return err
	}
	if book == nil {
		return errors.New("buku tidak ditemukan")
	}
	if book.Stock <= 0 {
		return errors.New("stok buku habis, peminjaman ditolak") // F-LOAN-03
	}

	// 2. Catat peminjaman (borrow_date & due_date diisi oleh repository)
	if _, err := s.loanRepo.Create(visitorID, staffID, bookID); err != nil {
		return err
	}

	// 3. Kurangi stok buku
	if err := s.bookRepo.DecrementStock(bookID); err != nil { // F-LOAN-04
		return err
	}

	// 4. Audit trail
	s.activityRepo.Log("Borrow Book", fmt.Sprintf(
		"Buku '%s' (ID:%d) dipinjam oleh visitor ID:%d", book.Title, bookID, visitorID,
	))

	return nil
}

// Return memproses pengembalian buku dan membuat invoice denda jika terlambat.
// F-RET-01: isi return_date, ubah status → 'returned'
// F-RET-02: tambah stok buku
// F-RET-03: buat invoice otomatis jika terlambat
func (s *LoanService) Return(loanID int) error {
	// 1. Ambil data loan
	loan, err := s.loanRepo.FindByID(loanID)
	if err != nil {
		return err
	}
	if loan == nil {
		return errors.New("data peminjaman tidak ditemukan")
	}
	if loan.Status == "returned" {
		return errors.New("buku sudah dikembalikan sebelumnya")
	}

	// 2. Tandai sebagai returned
	if err := s.loanRepo.MarkReturned(loanID); err != nil { // F-RET-01
		return err
	}

	// 3. Tambah stok buku
	if err := s.bookRepo.IncrementStock(loan.BookID); err != nil { // F-RET-02
		return err
	}

	// 4. Cek keterlambatan & buat invoice otomatis jika perlu
	returnDate := time.Now()
	if returnDate.After(loan.DueDate) { // F-RET-03
		hariTerlambat := int(returnDate.Sub(loan.DueDate).Hours() / 24)
		if hariTerlambat < 1 {
			hariTerlambat = 1
		}
		denda := float64(hariTerlambat) * DendaPerHari

		if _, err := s.invoiceRepo.Create(loan.VisitorID, loanID, denda); err != nil {
			return fmt.Errorf("buku berhasil dikembalikan, tapi gagal membuat invoice: %w", err)
		}

		s.activityRepo.Log("Return Book", fmt.Sprintf(
			"Buku '%s' dikembalikan terlambat %d hari, denda Rp%.0f",
			loan.Book.Title, hariTerlambat, denda,
		))
		return fmt.Errorf("buku dikembalikan terlambat %d hari — denda Rp%.0f telah dibuat", hariTerlambat, denda)
	}

	// Tepat waktu
	s.activityRepo.Log("Return Book", fmt.Sprintf(
		"Buku '%s' (loan ID:%d) dikembalikan tepat waktu", loan.Book.Title, loanID,
	))
	return nil
}

// GetActiveLoans mengambil semua peminjaman aktif (dipakai staff).
func (s *LoanService) GetActiveLoans() ([]interface{}, error) {
	return nil, nil // akan dipakai controller, di-wrap oleh FindAllActive
}
