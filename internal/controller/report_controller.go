package controller

import (
	"context"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type ReportController struct {
	reportService *service.ReportService
}

func NewReportController(reportService *service.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

func (c *ReportController) GetMostBorrowedBooks(ctx context.Context) ([]repository.MostBorrowedBook, error) {
	return c.reportService.GetMostBorrowedBooks(ctx)
}

func (c *ReportController) GetTotalPaidFines(ctx context.Context) (float64, error) {
	return c.reportService.GetTotalPaidFines(ctx)
}

func (c *ReportController) GetTopUsersByFines(ctx context.Context) ([]repository.UserFine, error) {
	return c.reportService.GetTopUsersByFines(ctx)
}
