-- ------------------------------------------------------------
-- 1. users
-- ------------------------------------------------------------
CREATE TABLE users (
    id       INT          NOT NULL AUTO_INCREMENT,
    name     VARCHAR(100) NOT NULL,
    email    VARCHAR(150) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role     VARCHAR(20)  NOT NULL,
    PRIMARY KEY (id)
);

-- ------------------------------------------------------------
-- 2. authors  (F-CAT-01)
-- ------------------------------------------------------------
CREATE TABLE authors (
    id          INT          NOT NULL AUTO_INCREMENT,
    name        VARCHAR(150) NOT NULL,
    birth_date  DATE,
    nationality VARCHAR(100),
    bio         TEXT,
    PRIMARY KEY (id)
);

-- ------------------------------------------------------------
-- 3. books  (F-CAT-02, F-CAT-03)
-- ------------------------------------------------------------
CREATE TABLE books (
    id        INT          NOT NULL AUTO_INCREMENT,
    isbn      VARCHAR(20)  NOT NULL UNIQUE,
    title     VARCHAR(255) NOT NULL,
    author_id INT          NOT NULL,
    genre     VARCHAR(100) NOT NULL,
    stock     INT          NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    CONSTRAINT fk_books_author
        FOREIGN KEY (author_id) REFERENCES authors(id)
        ON DELETE RESTRICT
);

-- ------------------------------------------------------------
-- 4. loans  (F-LOAN-01, F-LOAN-02)
-- ------------------------------------------------------------
CREATE TABLE loans (
    id          INT         NOT NULL AUTO_INCREMENT,
    visitor_id  INT         NOT NULL,
    staff_id    INT         NOT NULL,
    book_id     INT         NOT NULL,
    borrow_date DATE        NOT NULL,
    due_date    DATE        NOT NULL,
    return_date DATE,
    status      VARCHAR(20) NOT NULL,
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
-- 5. invoices  (F-RET-03, F-PAY-01)
-- ------------------------------------------------------------
CREATE TABLE invoices (
    id         INT           NOT NULL AUTO_INCREMENT,
    user_id    INT           NOT NULL,
    loan_id    INT           NOT NULL UNIQUE,
    amount     DECIMAL(10,2) NOT NULL,
    issue_date DATE          NOT NULL,
    status     VARCHAR(20),
    PRIMARY KEY (id),
    CONSTRAINT fk_invoices_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_invoices_loan
        FOREIGN KEY (loan_id) REFERENCES loans(id)
        ON DELETE RESTRICT
);

-- ------------------------------------------------------------
-- 6. payments  (F-PAY-02, F-PAY-03)
-- ------------------------------------------------------------
CREATE TABLE payments (
    id         INT           NOT NULL AUTO_INCREMENT,
    invoice_id INT           NOT NULL,
    amount     DECIMAL(10,2) NOT NULL,
    date       DATE          NOT NULL,
    method     VARCHAR(50)   NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_payments_invoice
        FOREIGN KEY (invoice_id) REFERENCES invoices(id)
        ON DELETE RESTRICT
);

-- ------------------------------------------------------------
-- 7. activity_logs  (F-LOG-01, F-LOG-02)
-- ------------------------------------------------------------
CREATE TABLE activity_logs (
    id          INT          NOT NULL AUTO_INCREMENT,
    title       VARCHAR(100) NOT NULL,
    description TEXT         NOT NULL,
    date        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
