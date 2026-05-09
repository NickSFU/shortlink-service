package main

import (
	"log"

	"github.com/NickSFU/shortlink-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
}
