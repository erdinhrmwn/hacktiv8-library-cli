package main

import (
	"fmt"
	"strings"

	"github.com/erdinhrmwn/hacktiv8-library-cli/config"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/database"
)

func main() {
	cfg := config.InitConfig()

	db, err := database.InitializeDB(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("Failed to ping database: %v\n", err)
		return
	}

	fmt.Println(strings.Repeat("=", 25))
	fmt.Println("📚 Hacktiv8 Library CLI")
	fmt.Println(strings.Repeat("=", 25))
}
