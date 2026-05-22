package repository

import (
	"database/sql"
	"time"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/model"
)

type LoanRepository struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) *LoanRepository {
	return &LoanRepository{db: db}
}

// Create mencatat peminjaman baru dengan borrow_date hari ini dan due_date +7 hari.
func (r *LoanRepository) Create(visitorID, staffID, bookID int) (int64, error) {
	borrowDate := time.Now()
	dueDate := borrowDate.AddDate(0, 0, 7)

	result, err := r.db.Exec(`
		INSERT INTO loans (visitor_id, staff_id, book_id, borrow_date, due_date, status)
		VALUES (?, ?, ?, ?, ?, 'active')`,
		visitorID, staffID, bookID,
		borrowDate.Format("2006-01-02"),
		dueDate.Format("2006-01-02"),
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// FindByID mengambil satu loan berdasarkan ID beserta relasi visitor dan buku.
func (r *LoanRepository) FindByID(id int) (*model.Loan, error) {
	var l model.Loan
	var visitor model.User
	var book model.Book
	var returnDate sql.NullTime

	err := r.db.QueryRow(`
		SELECT
			l.id, l.visitor_id, l.staff_id, l.book_id,
			l.borrow_date, l.due_date, l.return_date, l.status,
			u.name, b.title
		FROM loans l
		JOIN users u ON u.id = l.visitor_id
		JOIN books b ON b.id = l.book_id
		WHERE l.id = ?`, id,
	).Scan(
		&l.ID, &l.VisitorID, &l.StaffID, &l.BookID,
		&l.BorrowDate, &l.DueDate, &returnDate, &l.Status,
		&visitor.Name, &book.Title,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if returnDate.Valid {
		l.ReturnDate = &returnDate.Time
	}
	l.Visitor = &visitor
	l.Book = &book
	return &l, nil
}

// FindActiveByVisitorID mengambil semua loan aktif milik visitor tertentu.
func (r *LoanRepository) FindActiveByVisitorID(visitorID int) ([]model.Loan, error) {
	rows, err := r.db.Query(`
		SELECT
			l.id, l.visitor_id, l.staff_id, l.book_id,
			l.borrow_date, l.due_date, l.status,
			b.title
		FROM loans l
		JOIN books b ON b.id = l.book_id
		WHERE l.visitor_id = ? AND l.status = 'active'`, visitorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []model.Loan
	for rows.Next() {
		var l model.Loan
		var book model.Book
		if err := rows.Scan(
			&l.ID, &l.VisitorID, &l.StaffID, &l.BookID,
			&l.BorrowDate, &l.DueDate, &l.Status,
			&book.Title,
		); err != nil {
			return nil, err
		}
		l.Book = &book
		loans = append(loans, l)
	}
	return loans, nil
}

// FindAllActive mengambil semua loan aktif (dipakai staff untuk proses pengembalian).
func (r *LoanRepository) FindAllActive() ([]model.Loan, error) {
	rows, err := r.db.Query(`
		SELECT
			l.id, l.visitor_id, l.staff_id, l.book_id,
			l.borrow_date, l.due_date, l.status,
			u.name, b.title
		FROM loans l
		JOIN users u ON u.id = l.visitor_id
		JOIN books b ON b.id = l.book_id
		WHERE l.status = 'active'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []model.Loan
	for rows.Next() {
		var l model.Loan
		var visitor model.User
		var book model.Book
		if err := rows.Scan(
			&l.ID, &l.VisitorID, &l.StaffID, &l.BookID,
			&l.BorrowDate, &l.DueDate, &l.Status,
			&visitor.Name, &book.Title,
		); err != nil {
			return nil, err
		}
		l.Visitor = &visitor
		l.Book = &book
		loans = append(loans, l)
	}
	return loans, nil
}

// MarkReturned mengisi return_date dan mengubah status menjadi 'returned'.
func (r *LoanRepository) MarkReturned(loanID int) error {
	_, err := r.db.Exec(`
		UPDATE loans SET return_date = ?, status = 'returned' WHERE id = ?`,
		time.Now().Format("2006-01-02"), loanID,
	)
	return err
}
