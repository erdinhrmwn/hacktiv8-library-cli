package service

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

type BookService struct {
	bookRepository   *repository.BookRepository
	authorRepository *repository.AuthorRepository
}

func NewBookService(bookRepository *repository.BookRepository, authorRepository *repository.AuthorRepository) *BookService {
	return &BookService{bookRepository: bookRepository, authorRepository: authorRepository}
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]model.Book, error) {
	return s.bookRepository.GetAllBooks(ctx)
}

func (s *BookService) GetBookByID(ctx context.Context, id int) (*model.Book, error) {
	book, err := s.bookRepository.GetBookByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if book.ID == 0 {
		return nil, fmt.Errorf("buku tidak ditemukan")
	}
	return book, nil
}

func (s *BookService) SearchBook(ctx context.Context, query string) ([]model.Book, error) {
	return s.bookRepository.SearchBook(ctx, query)
}

func (s *BookService) AddBook(ctx context.Context, book *model.Book) error {
	if book.ISBN == "" {
		return fmt.Errorf("ISBN tidak boleh kosong")
	}
	if book.Title == "" {
		return fmt.Errorf("judul buku tidak boleh kosong")
	}

	existingISBN, err := s.bookRepository.FindByISBN(ctx, book.ISBN)
	if err != nil {
		return err
	}
	if existingISBN.ID != 0 {
		return fmt.Errorf("ISBN sudah terdaftar")
	}

	existing, err := s.authorRepository.GetAuthorByID(ctx, book.AuthorID)
	if err != nil {
		return err
	}
	if existing.ID == 0 {
		return fmt.Errorf("penulis tidak ditemukan")
	}

	return s.bookRepository.AddBook(ctx, book)
}

func (s *BookService) UpdateBook(ctx context.Context, book *model.Book) error {
	if book.Title == "" {
		return fmt.Errorf("judul buku tidak boleh kosong")
	}

	existing, err := s.bookRepository.GetBookByID(ctx, book.ID)
	if err != nil {
		return err
	}
	if existing.ID == 0 {
		return fmt.Errorf("buku tidak ditemukan")
	}

	return s.bookRepository.UpdateBook(ctx, book)
}

func (s *BookService) DeleteBook(ctx context.Context, id int) error {
	existing, err := s.bookRepository.GetBookByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.ID == 0 {
		return fmt.Errorf("buku tidak ditemukan")
	}

	return s.bookRepository.DeleteBook(ctx, id)
}
