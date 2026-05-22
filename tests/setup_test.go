package tests

import (
	"database/sql"
	"os"
	"testing"

	"github.com/erdinhrmwn/hacktiv8-library-cli/config"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/database"
	"github.com/erdinhrmwn/hacktiv8-library-cli/utils"
	_ "github.com/go-sql-driver/mysql"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	cfg := config.Config{
		DBHost:     "127.0.0.1",
		DBPort:     3306,
		DBUser:     "root",
		DBPassword: "",
		DBName:     "library_test",
	}

	var err error
	testDB, err = database.InitializeDB(cfg)
	if err != nil {
		panic("gagal konek database test: " + err.Error())
	}

	setupSchema(testDB)
	seedData(testDB)

	code := m.Run()

	cleanup(testDB)
	testDB.Close()
	os.Exit(code)
}

func setupSchema(db *sql.DB) {
	db.Exec("DROP DATABASE IF EXISTS library_test")
	db.Exec("CREATE DATABASE library_test")
	db.Exec("USE library_test")

	schema, _ := os.ReadFile("./sql/schema.sql")
	_, err := db.Exec(string(schema))
	if err != nil {
		panic("gagal setup schema: " + err.Error())
	}
}

func seedData(db *sql.DB) {
	hash, _ := utils.HashPassword("password")

	db.Exec(`INSERT INTO users (name, email, password, role) VALUES
		('Test Staff', 'staff@test.com', ?, 'staff'),
		('Test Visitor', 'visitor@test.com', ?, 'visitor')`, hash, hash)

	db.Exec(`INSERT INTO authors (name, birth_date, nationality, bio) VALUES
		('Author One', '1980-01-01', 'Indonesian', 'Bio one'),
		('Author Two', '1990-06-15', 'American', 'Bio two')`)

	db.Exec(`INSERT INTO books (isbn, title, author_id, genre, stock) VALUES
		('111-111', 'Book One', 1, 'Fantasy', 3),
		('222-222', 'Book Two', 2, 'Sci-Fi', 0),
		('333-333', 'Book Three', 1, 'Romance', 5)`)
}

func cleanup(db *sql.DB) {
	db.Exec("DROP DATABASE IF EXISTS library_test")
}
