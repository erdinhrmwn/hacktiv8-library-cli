-- ============================================================
--  Seed Data
-- ============================================================

-- Password untuk semua user: 'password'

INSERT INTO users (name, email, password, role) VALUES
('Admin Staff', 'admin@library.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'staff'),
('John Doe', 'visitor@library.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'visitor'),
('Jane Doe', 'jane@library.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'visitor'),
('Bob Smith', 'bob@library.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'visitor'),
('Alice Brown', 'alice@library.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'visitor'),
('Charlie Kim', 'charlie@library.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'visitor');

INSERT INTO authors (name, birth_date, nationality, bio) VALUES
('J.K. Rowling', '1965-07-31', 'British', 'Penulis seri Harry Potter'),
('George Orwell', '1903-06-25', 'British', 'Penulis 1984 dan Animal Farm'),
('J.R.R. Tolkien', '1892-01-03', 'British', 'Penulis The Lord of the Rings'),
('Andrea Hirata', '1967-10-24', 'Indonesia', 'Penulis Laskar Pelangi'),
('Pramoedya Ananta Toer', '1925-02-06', 'Indonesia', 'Penulis Bumi Manusia');

INSERT INTO books (isbn, title, author_id, genre, stock) VALUES
('978-0747532743', 'Harry Potter and the Philosopher''s Stone', 1, 'Fantasy', 5),
('978-0747538499', 'Harry Potter and the Chamber of Secrets', 1, 'Fantasy', 3),
('978-0451524935', '1984', 2, 'Sci-Fi', 4),
('978-0451526342', 'Animal Farm', 2, 'Sci-Fi', 6),
('978-0547928227', 'The Fellowship of the Ring', 3, 'Fantasy', 3),
('978-0547928203', 'The Hobbit', 3, 'Fantasy', 5),
('978-9793062792', 'Laskar Pelangi', 4, 'Slice of Life', 4),
('978-9799731234', 'Bumi Manusia', 5, 'Romance', 3);

-- Loans: berbagai case (on-time, terlambat, active, overdue)
-- Staff id=1, Visitors: John(2), Jane(3), Bob(4), Alice(5), Charlie(6)

INSERT INTO loans (visitor_id, staff_id, book_id, borrow_date, due_date, return_date, status) VALUES
-- John: Harry Potter 1 - returned on time
(2, 1, 1, DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), 'returned'),
-- Jane: 1984 - returned late (3 days)
(3, 1, 3, DATE_SUB(NOW(), INTERVAL 14 DAY), DATE_SUB(NOW(), INTERVAL 7 DAY), DATE_SUB(NOW(), INTERVAL 4 DAY), 'returned'),
-- John: Animal Farm - returned late (5 days)
(2, 1, 4, DATE_SUB(NOW(), INTERVAL 15 DAY), DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), 'returned'),
-- John: Harry Potter 1 (again) - returned late (2 days)
(2, 1, 1, DATE_SUB(NOW(), INTERVAL 12 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), 'returned'),
-- Jane: Harry Potter 1 - returned late (1 day)
(3, 1, 1, DATE_SUB(NOW(), INTERVAL 9 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 'returned'),
-- Bob: Harry Potter 2 - returned on time
(4, 1, 2, DATE_SUB(NOW(), INTERVAL 8 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 'returned'),
-- Alice: 1984 - returned late (4 days)
(5, 1, 3, DATE_SUB(NOW(), INTERVAL 16 DAY), DATE_SUB(NOW(), INTERVAL 9 DAY), DATE_SUB(NOW(), INTERVAL 5 DAY), 'returned'),
-- Charlie: Laskar Pelangi - returned late (2 days)
(6, 1, 7, DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 1 DAY), 'returned'),
-- Jane: The Hobbit - active (not overdue yet)
(3, 1, 6, DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_ADD(NOW(), INTERVAL 4 DAY), NULL, 'active'),
-- John: LOTR - active (not overdue)
(2, 1, 5, DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_ADD(NOW(), INTERVAL 5 DAY), NULL, 'active'),
-- Bob: Animal Farm - active (not overdue)
(4, 1, 4, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 6 DAY), NULL, 'active'),
-- Alice: The Hobbit - active but OVERDUE (3 days)
(5, 1, 6, DATE_SUB(NOW(), INTERVAL 10 DAY), DATE_SUB(NOW(), INTERVAL 3 DAY), NULL, 'active'),
-- Charlie: Bumi Manusia - active (not overdue)
(6, 1, 8, DATE_SUB(NOW(), INTERVAL 2 DAY), DATE_ADD(NOW(), INTERVAL 5 DAY), NULL, 'active');

-- Invoices for late returns

INSERT INTO invoices (user_id, loan_id, amount, issue_date, status) VALUES
-- Jane: 1984 late 3 days = Rp15,000 (paid)
(3, 2, 15000, DATE_SUB(NOW(), INTERVAL 4 DAY), 'paid'),
-- John: Animal Farm late 5 days = Rp25,000 (unpaid)
(2, 3, 25000, DATE_SUB(NOW(), INTERVAL 5 DAY), 'unpaid'),
-- John: Harry Potter 1 (2nd) late 2 days = Rp10,000 (paid)
(2, 4, 10000, DATE_SUB(NOW(), INTERVAL 3 DAY), 'paid'),
-- Jane: Harry Potter 1 late 1 day = Rp5,000 (paid)
(3, 5, 5000, DATE_SUB(NOW(), INTERVAL 1 DAY), 'paid'),
-- Alice: 1984 late 4 days = Rp20,000 (paid)
(5, 7, 20000, DATE_SUB(NOW(), INTERVAL 5 DAY), 'paid'),
-- Charlie: Laskar Pelangi late 2 days = Rp10,000 (unpaid)
(6, 8, 10000, DATE_SUB(NOW(), INTERVAL 1 DAY), 'unpaid');

-- Payments for paid invoices

INSERT INTO payments (invoice_id, amount, date, method) VALUES
(1, 15000, DATE_SUB(NOW(), INTERVAL 4 DAY), 'cash'),
(3, 10000, DATE_SUB(NOW(), INTERVAL 3 DAY), 'transfer'),
(4, 5000, DATE_SUB(NOW(), INTERVAL 1 DAY), 'cash'),
(5, 20000, DATE_SUB(NOW(), INTERVAL 5 DAY), 'transfer');
