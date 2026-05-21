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

func (c *BookController) AddBook(ctx context.Context, isbn, title, genre string, authorID, stock int) error {
	if strings.TrimSpace(isbn) == "" {
		return fmt.Errorf("ISBN tidak boleh kosong")
	}
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("judul buku tidak boleh kosong")
	}
	if authorID <= 0 {
		return fmt.Errorf("ID penulis tidak valid")
	}
	if stock < 0 {
		return fmt.Errorf("stok tidak boleh negatif")
	}

	book := &model.Book{
		ISBN:     strings.TrimSpace(isbn),
		Title:    strings.TrimSpace(title),
		AuthorID: authorID,
		Genre:    genre,
		Stock:    stock,
	}
	return c.bookService.AddBook(ctx, book)
}

func (c *BookController) UpdateBook(ctx context.Context, id int, isbn, title, genre string, authorID, stock int) error {
	if id <= 0 {
		return fmt.Errorf("ID buku tidak valid")
	}
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("judul buku tidak boleh kosong")
	}
	if authorID <= 0 {
		return fmt.Errorf("ID penulis tidak valid")
	}
	if stock < 0 {
		return fmt.Errorf("stok tidak boleh negatif")
	}

	book := &model.Book{
		ID:       id,
		ISBN:     strings.TrimSpace(isbn),
		Title:    strings.TrimSpace(title),
		AuthorID: authorID,
		Genre:    genre,
		Stock:    stock,
	}
	return c.bookService.UpdateBook(ctx, book)
}

func (c *BookController) DeleteBook(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("ID buku tidak valid")
	}
	return c.bookService.DeleteBook(ctx, id)
}
