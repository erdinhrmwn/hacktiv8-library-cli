package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/erdinhrmwn/hacktiv8-library-cli/internal/container"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	app := container.New()

	fmt.Println("📚 Hacktiv8 Library CLI — v0.1.0")
	app.Run(ctx)
}
