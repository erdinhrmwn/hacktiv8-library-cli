package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type BookController struct {
	bookService *service.BookService
}

func NewBookController(bookService *service.BookService) *BookController {
	return &BookController{bookService: bookService}
}

func (c *BookController) GetAllBooks(ctx context.Context) ([]model.Book, error) {
	return c.bookService.GetAllBooks(ctx)
}

func (c *BookController) GetBookByID(ctx context.Context, id int) (*model.Book, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID buku tidak valid")
	}
	return c.bookService.GetBookByID(ctx, id)
}

func (c *BookController) SearchBook(ctx context.Context, query string) ([]model.Book, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("kata kunci pencarian tidak boleh kosong")
	}
	return c.bookService.SearchBook(ctx, query)
}

type AddBookInput struct {
	ISBN     string
	Title    string
	Genre    string
	AuthorID int
	Stock    int
}

func (c *BookController) AddBook(ctx context.Context, input AddBookInput) error {
	if strings.TrimSpace(input.ISBN) == "" {
		return fmt.Errorf("ISBN tidak boleh kosong")
	}
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("judul buku tidak boleh kosong")
	}
	if input.AuthorID <= 0 {
		return fmt.Errorf("ID penulis tidak valid")
	}
	if input.Stock < 0 {
		return fmt.Errorf("stok tidak boleh negatif")
	}

	book := &model.Book{
		ISBN:     strings.TrimSpace(input.ISBN),
		Title:    strings.TrimSpace(input.Title),
		AuthorID: input.AuthorID,
		Genre:    input.Genre,
		Stock:    input.Stock,
	}
	return c.bookService.AddBook(ctx, book)
}

type UpdateBookInput struct {
	ID       int
	ISBN     string
	Title    string
	Genre    string
	AuthorID int
	Stock    int
}

func (c *BookController) UpdateBook(ctx context.Context, input UpdateBookInput) error {
	if input.ID <= 0 {
		return fmt.Errorf("ID buku tidak valid")
	}
	if strings.TrimSpace(input.Title) == "" {
		return fmt.Errorf("judul buku tidak boleh kosong")
	}
	if input.AuthorID <= 0 {
		return fmt.Errorf("ID penulis tidak valid")
	}
	if input.Stock < 0 {
		return fmt.Errorf("stok tidak boleh negatif")
	}

	book := &model.Book{
		ID:       input.ID,
		ISBN:     strings.TrimSpace(input.ISBN),
		Title:    strings.TrimSpace(input.Title),
		AuthorID: input.AuthorID,
		Genre:    input.Genre,
		Stock:    input.Stock,
	}
	return c.bookService.UpdateBook(ctx, book)
}

func (c *BookController) DeleteBook(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("ID buku tidak valid")
	}
	return c.bookService.DeleteBook(ctx, id)
}
