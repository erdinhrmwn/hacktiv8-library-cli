package service

import (
	"context"
	"fmt"
	"time"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

const FinePerDay = 5000.0

type LoanService struct {
	loanRepository    *repository.LoanRepository
	bookRepository    *repository.BookRepository
	invoiceRepository *repository.InvoiceRepository
}

func NewLoanService(
	loanRepository *repository.LoanRepository,
	bookRepository *repository.BookRepository,
	invoiceRepository *repository.InvoiceRepository,
) *LoanService {
	return &LoanService{
		loanRepository:    loanRepository,
		bookRepository:    bookRepository,
		invoiceRepository: invoiceRepository,
	}
}

func (s *LoanService) Borrow(ctx context.Context, visitorID, staffID, bookID int) error {
	book, err := s.bookRepository.GetBookByID(ctx, bookID)
	if err != nil {
		return err
	}
	if book.ID == 0 {
		return fmt.Errorf("buku tidak ditemukan")
	}
	if book.Stock <= 0 {
		return fmt.Errorf("stok buku habis, peminjaman ditolak")
	}

	if err := s.loanRepository.CreateLoan(ctx, visitorID, staffID, bookID); err != nil {
		return err
	}
	if err := s.bookRepository.DecrementStock(ctx, bookID); err != nil {
		return err
	}

	return nil
}

func (s *LoanService) Return(ctx context.Context, loanID int) (*float64, error) {
	loan, err := s.loanRepository.GetLoanByID(ctx, loanID)
	if err != nil {
		return nil, err
	}
	if loan.ID == 0 {
		return nil, fmt.Errorf("data peminjaman tidak ditemukan")
	}
	if loan.Status == "returned" {
		return nil, fmt.Errorf("buku sudah dikembalikan")
	}

	if err := s.loanRepository.ReturnLoan(ctx, loanID); err != nil {
		return nil, err
	}
	if err := s.bookRepository.IncrementStock(ctx, loan.BookID); err != nil {
		return nil, err
	}

	now := time.Now()
	if now.After(loan.DueDate) {
		daysLate := max(int(now.Sub(loan.DueDate).Hours()/24), 1)
		fine := float64(daysLate) * FinePerDay

		if err := s.invoiceRepository.CreateInvoice(ctx, loan.VisitorID, loanID, fine); err != nil {
			return nil, fmt.Errorf("buku dikembalikan, tapi gagal membuat invoice: %w", err)
		}

		return &fine, fmt.Errorf("buku terlambat %d hari — denda Rp%.0f", daysLate, fine)
	}

	return nil, nil
}

func (s *LoanService) GetActiveByVisitorID(ctx context.Context, visitorID int) ([]model.Loan, error) {
	loans, err := s.loanRepository.GetActiveLoansByVisitorID(ctx, visitorID)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (s *LoanService) GetLoanByID(ctx context.Context, id int) (*model.Loan, error) {
	loan, err := s.loanRepository.GetLoanByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if loan.ID == 0 {
		return nil, fmt.Errorf("loan tidak ditemukan")
	}
	return loan, nil
}

func (s *LoanService) GetAllLoans(ctx context.Context) ([]model.Loan, error) {
	loans, err := s.loanRepository.GetAllLoans(ctx)
	if err != nil {
		return nil, err
	}
	return loans, nil
}
