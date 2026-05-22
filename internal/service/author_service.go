package service

import (
	"context"
	"fmt"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/repository"
)

type AuthorService struct {
	authorRepository *repository.AuthorRepository
}

func NewAuthorService(authorRepository *repository.AuthorRepository) *AuthorService {
	return &AuthorService{authorRepository: authorRepository}
}

func (s *AuthorService) GetAllAuthors(ctx context.Context) ([]model.Author, error) {
	authors, err := s.authorRepository.GetAllAuthors(ctx)
	if err != nil {
		return nil, err
	}
	return authors, nil
}

func (s *AuthorService) GetAuthorByID(ctx context.Context, id int) (*model.Author, error) {
	author, err := s.authorRepository.GetAuthorByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if author.ID == 0 {
		return nil, fmt.Errorf("penulis tidak ditemukan")
	}
	return author, nil
}

func (s *AuthorService) SearchAuthor(ctx context.Context, query string) ([]model.Author, error) {
	authors, err := s.authorRepository.SearchAuthor(ctx, query)
	if err != nil {
		return nil, err
	}
	return authors, nil
}

func (s *AuthorService) AddAuthor(ctx context.Context, author *model.Author) error {
	if author.Name == "" {
		return fmt.Errorf("nama penulis tidak boleh kosong")
	}

	return s.authorRepository.AddAuthor(ctx, author)
}

func (s *AuthorService) UpdateAuthor(ctx context.Context, author *model.Author) error {
	if author.Name == "" {
		return fmt.Errorf("nama penulis tidak boleh kosong")
	}

	existing, err := s.authorRepository.GetAuthorByID(ctx, author.ID)
	if err != nil {
		return err
	}
	if existing.ID == 0 {
		return fmt.Errorf("penulis tidak ditemukan")
	}

	return s.authorRepository.UpdateAuthor(ctx, author)
}

func (s *AuthorService) DeleteAuthor(ctx context.Context, id int) error {
	existing, err := s.authorRepository.GetAuthorByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.ID == 0 {
		return fmt.Errorf("penulis tidak ditemukan")
	}

	return s.authorRepository.DeleteAuthor(ctx, id)
}

func (s *AuthorService) GetAuthorBooks(ctx context.Context, authorID int) ([]model.Book, error) {
	books, err := s.authorRepository.GetAuthorBooks(ctx, authorID)
	if err != nil {
		return nil, err
	}
	return books, nil
}
