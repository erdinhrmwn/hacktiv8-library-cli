package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
	"github.com/olekukonko/tablewriter"
)

type LoanController struct {
	loanService    *service.LoanService
	loanRepository *repository.LoanRepository
}

func NewLoanController(loanService *service.LoanService, loanRepository *repository.LoanRepository) *LoanController {
	return &LoanController{loanService: loanService, loanRepository: loanRepository}
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

func (c *LoanController) GetActiveByVisitorID(ctx context.Context, visitorID int) ([]model.Loan, error) {
	return c.loanService.GetActiveByVisitorID(ctx, visitorID)
}

func (c *LoanController) GetAll(ctx context.Context) ([]model.Loan, error) {
	loans, err := c.loanRepository.GetAllLoans(ctx)
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func RenderLoanTable(data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"ID", "Visitor", "Book", "Borrow", "Due", "Status"})
	for _, row := range data {
		table.Append([]string{row[0], row[1], row[2], row[3], row[4], row[5]})
	}
	table.Render()
}
