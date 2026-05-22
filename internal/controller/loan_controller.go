package controller

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type LoanController struct {
	loanService *service.LoanService
}

func NewLoanController(loanService *service.LoanService) *LoanController {
	return &LoanController{loanService: loanService}
}

func (c *LoanController) BorrowBook(ctx context.Context, visitorID, staffID, bookID int) error {
	if visitorID <= 0 {
		return fmt.Errorf("visitor tidak valid")
	}
	if staffID <= 0 {
		return fmt.Errorf("staff tidak valid")
	}
	if bookID <= 0 {
		return fmt.Errorf("buku tidak valid")
	}
	return c.loanService.Borrow(ctx, visitorID, staffID, bookID)
}

func (c *LoanController) ReturnBook(ctx context.Context, loanID int) (*float64, error) {
	if loanID <= 0 {
		return nil, fmt.Errorf("loan ID tidak valid")
	}
	return c.loanService.Return(ctx, loanID)
}

func (c *LoanController) GetLoansByVisitorID(ctx context.Context, visitorID int) ([]model.Loan, error) {
	if visitorID <= 0 {
		return nil, fmt.Errorf("visitor ID tidak valid")
	}
	return c.loanService.GetLoansByVisitorID(ctx, visitorID)
}

func (c *LoanController) GetActiveByVisitorID(ctx context.Context, visitorID int) ([]model.Loan, error) {
	return c.loanService.GetActiveByVisitorID(ctx, visitorID)
}

func (c *LoanController) GetLoanByID(ctx context.Context, id int) (*model.Loan, error) {
	if id <= 0 {
		return nil, fmt.Errorf("loan ID tidak valid")
	}
	return c.loanService.GetLoanByID(ctx, id)
}

func (c *LoanController) GetAll(ctx context.Context) ([]model.Loan, error) {
	return c.loanService.GetAllLoans(ctx)
}
