package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/erdinhrmwn/hacktiv8-library-cli/config"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/container"
	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/database"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.InitConfig()

	db, err := database.InitializeDB(cfg)
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke database: %v", err)
	}
	defer db.Close()

	fmt.Println()
	fmt.Println(strings.Repeat("═", 50))
	fmt.Println("   📚  HACKTIV8 LIBRARY CLI  —  v1.0.0")
	fmt.Println(strings.Repeat("═", 50))
	fmt.Println()

	app := container.New(db)
	app.Run(ctx)
}
