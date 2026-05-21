package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/service"
)

type AuthorController struct {
	authorService *service.AuthorService
}

func NewAuthorController(authorService *service.AuthorService) *AuthorController {
	return &AuthorController{authorService: authorService}
}

func (c *AuthorController) GetAllAuthors(ctx context.Context) ([]model.Author, error) {
	return c.authorService.GetAllAuthors(ctx)
}

func (c *AuthorController) GetAuthorByID(ctx context.Context, id int) (*model.Author, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID penulis tidak valid")
	}
	return c.authorService.GetAuthorByID(ctx, id)
}

func (c *AuthorController) SearchAuthor(ctx context.Context, query string) ([]model.Author, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("kata kunci pencarian tidak boleh kosong")
	}
	return c.authorService.SearchAuthor(ctx, query)
}

type AddAuthorInput struct {
	Name        string
	BirthDate   string
	Nationality string
	Bio         string
}

func (c *AuthorController) AddAuthor(ctx context.Context, input AddAuthorInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("nama penulis tidak boleh kosong")
	}

	author := &model.Author{
		Name:        strings.TrimSpace(input.Name),
		BirthDate:   input.BirthDate,
		Nationality: input.Nationality,
		Bio:         input.Bio,
	}
	return c.authorService.AddAuthor(ctx, author)
}

type UpdateAuthorInput struct {
	ID          int
	Name        string
	BirthDate   string
	Nationality string
	Bio         string
}

func (c *AuthorController) UpdateAuthor(ctx context.Context, input UpdateAuthorInput) error {
	if input.ID <= 0 {
		return fmt.Errorf("ID penulis tidak valid")
	}
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("nama penulis tidak boleh kosong")
	}

	author := &model.Author{
		ID:          input.ID,
		Name:        strings.TrimSpace(input.Name),
		BirthDate:   input.BirthDate,
		Nationality: input.Nationality,
		Bio:         input.Bio,
	}
	return c.authorService.UpdateAuthor(ctx, author)
}

func (c *AuthorController) DeleteAuthor(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("ID penulis tidak valid")
	}
	return c.authorService.DeleteAuthor(ctx, id)
}

func (c *AuthorController) GetAuthorBooks(ctx context.Context, authorID int) ([]model.Book, error) {
	if authorID <= 0 {
		return nil, fmt.Errorf("ID penulis tidak valid")
	}
	return c.authorService.GetAuthorBooks(ctx, authorID)
}
