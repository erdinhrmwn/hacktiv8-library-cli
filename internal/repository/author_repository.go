package repository

import (
	"context"
	"database/sql"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
)

type AuthorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) GetAllAuthors(ctx context.Context) ([]model.Author, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM authors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []model.Author
	for rows.Next() {
		var author model.Author
		if err := rows.Scan(&author.ID, &author.Name, &author.BirthDate, &author.Nationality, &author.Bio); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *AuthorRepository) GetAuthorByID(ctx context.Context, id int) (*model.Author, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM authors WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var author model.Author
	if rows.Next() {
		if err := rows.Scan(&author.ID, &author.Name, &author.BirthDate, &author.Nationality, &author.Bio); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &author, nil
}

func (r *AuthorRepository) SearchAuthor(ctx context.Context, query string) ([]model.Author, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM authors WHERE name LIKE ?`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []model.Author
	for rows.Next() {
		var a model.Author
		if err := rows.Scan(&a.ID, &a.Name, &a.BirthDate, &a.Nationality, &a.Bio); err != nil {
			return nil, err
		}
		authors = append(authors, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func (r *AuthorRepository) AddAuthor(ctx context.Context, author *model.Author) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO authors (name, birth_date, nationality, bio)
		VALUES (?, ?, ?, ?)`, author.Name, author.BirthDate, author.Nationality, author.Bio)
	return err
}

func (r *AuthorRepository) UpdateAuthor(ctx context.Context, author *model.Author) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE authors
		SET name = ?, birth_date = ?, nationality = ?, bio = ?
		WHERE id = ?`, author.Name, author.BirthDate, author.Nationality, author.Bio, author.ID)
	return err
}

func (r *AuthorRepository) DeleteAuthor(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM authors WHERE id = ?`, id)
	return err
}

func (r *AuthorRepository) GetAuthorBooks(ctx context.Context, authorID int) ([]model.Book, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM books WHERE author_id = ?`, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
		if err := rows.Scan(&book.ID, &book.ISBN, &book.Title, &book.AuthorID, &book.Genre, &book.Stock); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
