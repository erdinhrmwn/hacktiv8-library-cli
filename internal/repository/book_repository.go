package repository

import (
	"context"
	"database/sql"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetAllBooks(ctx context.Context) ([]model.Book, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM books")
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

func (r *BookRepository) GetBookByID(ctx context.Context, id int) (*model.Book, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			b.id, b.isbn, b.title, b.author_id, b.genre, b.stock,
			a.id as author_id, a.name as author_name, a.birth_date, a.nationality
		FROM books b
		JOIN authors a ON b.author_id = a.id
		WHERE b.id = ?`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var book model.Book
	if rows.Next() {
		var author model.Author
		if err := rows.Scan(
			&book.ID, &book.ISBN, &book.Title, &book.AuthorID, &book.Genre, &book.Stock,
			&author.ID, &author.Name, &author.BirthDate, &author.Nationality,
		); err != nil {
			return nil, err
		}

		book.Author = &author
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *BookRepository) SearchBook(ctx context.Context, query string) ([]model.Book, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM books WHERE title LIKE ?`, "%"+query+"%")
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

func (r *BookRepository) FindByISBN(ctx context.Context, isbn string) (*model.Book, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM books WHERE isbn = ?`, isbn)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var book model.Book
	if rows.Next() {
		if err := rows.Scan(&book.ID, &book.ISBN, &book.Title, &book.AuthorID, &book.Genre, &book.Stock); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *BookRepository) AddBook(ctx context.Context, book *model.Book) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO books (isbn, title, author_id, genre, stock)
		VALUES (?, ?, ?, ?, ?)`, book.ISBN, book.Title, book.AuthorID, book.Genre, book.Stock)
	return err
}

func (r *BookRepository) UpdateBook(ctx context.Context, book *model.Book) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE books
		SET isbn = ?, title = ?, author_id = ?, genre = ?, stock = ?
		WHERE id = ?`, book.ISBN, book.Title, book.AuthorID, book.Genre, book.Stock, book.ID)
	return err
}

func (r *BookRepository) DeleteBook(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM books WHERE id = ?`, id)
	return err
}
