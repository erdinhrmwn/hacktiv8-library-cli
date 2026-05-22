package service

import (
	"context"
	"fmt"
	"time"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

const DendaPerHari = 5000.0

type LoanService struct {
	loanRepo    *repository.LoanRepository
	bookRepo    *repository.BookRepository
	invoiceRepo *repository.InvoiceRepository
}

func NewLoanService(
	lr *repository.LoanRepository,
	br *repository.BookRepository,
	ir *repository.InvoiceRepository,
) *LoanService {
	return &LoanService{
		loanRepo:    lr,
		bookRepo:    br,
		invoiceRepo: ir,
	}
}

func (s *LoanService) Borrow(ctx context.Context, visitorID, staffID, bookID int) error {
	book, err := s.bookRepo.GetBookByID(ctx, bookID)
	if err != nil {
		return err
	}
	if book.ID == 0 {
		return fmt.Errorf("buku tidak ditemukan")
	}
	if book.Stock <= 0 {
		return fmt.Errorf("stok buku habis, peminjaman ditolak")
	}

	if err := s.loanRepo.Create(ctx, visitorID, staffID, bookID); err != nil {
		return err
	}
	if err := s.bookRepo.DecrementStock(ctx, bookID); err != nil {
		return err
	}

	return nil
}

func (s *LoanService) Return(ctx context.Context, loanID int) (*float64, error) {
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return nil, err
	}
	if loan.ID == 0 {
		return nil, fmt.Errorf("data peminjaman tidak ditemukan")
	}
	if loan.Status == "returned" {
		return nil, fmt.Errorf("buku sudah dikembalikan")
	}

	if err := s.loanRepo.Return(ctx, loanID); err != nil {
		return nil, err
	}
	if err := s.bookRepo.IncrementStock(ctx, loan.BookID); err != nil {
		return nil, err
	}

	now := time.Now()
	if now.After(loan.DueDate) {
		hariTerlambat := int(now.Sub(loan.DueDate).Hours() / 24)
		if hariTerlambat < 1 {
			hariTerlambat = 1
		}
		denda := float64(hariTerlambat) * DendaPerHari

		if err := s.invoiceRepo.Create(ctx, loan.VisitorID, loanID, denda); err != nil {
			return nil, fmt.Errorf("buku dikembalikan, tapi gagal membuat invoice: %w", err)
		}

		return &denda, fmt.Errorf("buku terlambat %d hari — denda Rp%.0f", hariTerlambat, denda)
	}

	return nil, nil
}

func (s *LoanService) GetActiveByVisitorID(ctx context.Context, visitorID int) ([]model.Loan, error) {
	loans, err := s.loanRepo.GetActiveByVisitorID(ctx, visitorID)
	if err != nil {
		return nil, err
	}
	return loans, nil
}
