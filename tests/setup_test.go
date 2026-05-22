package tests

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var testDB *sql.DB

const dsnRoot = "root@tcp(127.0.0.1:3306)/?parseTime=true"
const dsnTest = "root@tcp(127.0.0.1:3306)/library_test?parseTime=true&multiStatements=true"

func TestMain(m *testing.M) {
	db, err := sql.Open("mysql", dsnRoot)
	if err != nil {
		panic("gagal konek root: " + err.Error())
	}

	db.Exec("DROP DATABASE IF EXISTS library_test")
	db.Exec("CREATE DATABASE library_test")
	db.Close()

	testDB, err = sql.Open("mysql", dsnTest)
	if err != nil {
		panic("gagal konek test DB: " + err.Error())
	}

	setupSchema(testDB)
	seedData(testDB)

	code := m.Run()

	cleanup()
	testDB.Close()
	os.Exit(code)
}

func setupSchema(db *sql.DB) {
	schema, err := os.ReadFile("../sql/schema.sql")
	if err != nil {
		panic("gagal baca schema: " + err.Error())
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		panic("gagal setup schema: " + err.Error())
	}
}

func seedData(db *sql.DB) {
	db.Exec(`INSERT INTO users (name, email, password, role) VALUES
		('Test Staff', 'staff@test.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'staff'),
		('Test Visitor', 'visitor@test.com', '$2a$10$wT0s4lsAd/tN/LcFwDDLzuWk8aRTdZT1/M8nxSG2aIFVVv9eIVgWy', 'visitor')`)

	db.Exec(`INSERT INTO authors (name, birth_date, nationality, bio) VALUES
		('Author One', '1980-01-01', 'Indonesian', 'Bio one'),
		('Author Two', '1990-06-15', 'American', 'Bio two')`)

	db.Exec(`INSERT INTO books (isbn, title, author_id, genre, stock) VALUES
		('111-111', 'Book One', 1, 'Fantasy', 3),
		('222-222', 'Book Two', 2, 'Sci-Fi', 0),
		('333-333', 'Book Three', 1, 'Romance', 5)`)
}

func cleanup() {
	db, _ := sql.Open("mysql", dsnRoot)
	if db != nil {
		db.Exec("DROP DATABASE IF EXISTS library_test")
		db.Close()
	}
}
