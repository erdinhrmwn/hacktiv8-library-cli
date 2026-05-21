package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

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

	app := container.New(db)

	fmt.Println("📚 Hacktiv8 Library CLI — v0.1.0")
	app.Run(ctx)
}
