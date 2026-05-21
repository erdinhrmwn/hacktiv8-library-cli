
CREATE DATABASE library;

USE library;

-- ------------------------------------------------------------
-- 1. users
-- ------------------------------------------------------------
CREATE TABLE users (
    id         INT          NOT NULL AUTO_INCREMENT,
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(150) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,          -- bcrypt hash (F-AUTH-02)
    role       ENUM('staff','visitor') NOT NULL DEFAULT 'visitor',  -- F-AUTH-04
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                            ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- ------------------------------------------------------------
-- 2. authors  
-- ------------------------------------------------------------
CREATE TABLE authors (
    id         INT          NOT NULL AUTO_INCREMENT,
    name       VARCHAR(150) NOT NULL,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                            ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

-- ------------------------------------------------------------
-- 3. books 
-- ------------------------------------------------------------
CREATE TABLE books (
    id         INT          NOT NULL AUTO_INCREMENT,
    isbn       VARCHAR(20)  NOT NULL UNIQUE,   -- F-CAT-03: isbn harus unik
    title      VARCHAR(255) NOT NULL,
    author_id  INT          NOT NULL,           -- F-CAT-03: harus merujuk author valid
    stock      INT          NOT NULL DEFAULT 0 CHECK (stock >= 0),
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
                            ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_books_author
        FOREIGN KEY (author_id) REFERENCES authors(id)
        ON DELETE RESTRICT                      -- Aturan Bisnis: ON DELETE RESTRICT
);

-- ------------------------------------------------------------
-- 4. loans  
-- ------------------------------------------------------------
CREATE TABLE loans (
    id          INT      NOT NULL AUTO_INCREMENT,
    visitor_id  INT      NOT NULL,
    staff_id    INT      NOT NULL,
    book_id     INT      NOT NULL,
    borrow_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- F-LOAN-02
    due_date    DATETIME NOT NULL,             -- F-LOAN-02: borrow_date + 7 hari (di-set service)
    return_date DATETIME,                      -- F-RET-01: NULL = belum dikembalikan
    status      ENUM('active','returned') NOT NULL DEFAULT 'active',  -- F-RET-01
    PRIMARY KEY (id),
    CONSTRAINT fk_loans_visitor
        FOREIGN KEY (visitor_id) REFERENCES users(id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_loans_staff
        FOREIGN KEY (staff_id) REFERENCES users(id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_loans_book
        FOREIGN KEY (book_id) REFERENCES books(id)
        ON DELETE RESTRICT
);

-- ------------------------------------------------------------
-- 5. invoices  
-- ------------------------------------------------------------
CREATE TABLE invoices (
    id         INT           NOT NULL AUTO_INCREMENT,
    visitor_id INT           NOT NULL,
    loan_id    INT           NOT NULL UNIQUE,  -- 1 loan hanya bisa punya 1 invoice
    amount     DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    status     ENUM('unpaid','paid') NOT NULL DEFAULT 'unpaid',  -- F-PAY-03
    created_at DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_invoices_visitor
        FOREIGN KEY (visitor_id) REFERENCES users(id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_invoices_loan
        FOREIGN KEY (loan_id) REFERENCES loans(id)
        ON DELETE RESTRICT
);

-- ------------------------------------------------------------
-- 6. payments  
-- ------------------------------------------------------------
CREATE TABLE payments (
    id         INT           NOT NULL AUTO_INCREMENT,
    invoice_id INT           NOT NULL,
    amount_paid DECIMAL(10,2) NOT NULL,        -- F-PAY-02: amount_paid
    method     ENUM('cash','transfer') NOT NULL, -- F-PAY-02: method cash/transfer
    paid_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_payments_invoice
        FOREIGN KEY (invoice_id) REFERENCES invoices(id)
        ON DELETE RESTRICT
);

-- ------------------------------------------------------------
-- 7. activity_logs  
-- ------------------------------------------------------------
CREATE TABLE activity_logs (
    id          INT          NOT NULL AUTO_INCREMENT,
    `key`       VARCHAR(100) NOT NULL,   -- F-LOG-02: jenis aksi (Login, Add Book, dst)
    description TEXT         NOT NULL,   -- F-LOG-02: detail aksi
    date        DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- F-LOG-02
    PRIMARY KEY (id)
);

-- ============================================================
--  Seed Data
-- ============================================================


INSERT INTO users (name, email, password, role) VALUES
('Admin Staff', 'admin@library.com',
 '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lh2.',
 'staff');