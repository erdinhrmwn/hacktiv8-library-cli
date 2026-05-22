-- ============================================================
--  Seed Data
-- ============================================================

-- Password untuk semua user: 'password'

INSERT INTO users (name, email, password, role) VALUES
('Admin Staff', 'admin@library.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'staff'),
('John Doe', 'visitor@library.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'visitor'),
('Jane Doe', 'jane@library.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'visitor');

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
