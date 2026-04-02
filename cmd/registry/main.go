package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/stockyard-dev/stockyard-registry/internal/server"
	"github.com/stockyard-dev/stockyard-registry/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9220"
	}
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./registry-data"
	}

	db, err := store.Open(dataDir)
	if err != nil {
		log.Fatalf("registry: open database: %v", err)
	}
	defer db.Close()

	srv := server.New(db)

	fmt.Printf("\n  Registry — Self-hosted container and package registry\n")
	fmt.Printf("  ─────────────────────────────────\n")
	fmt.Printf("  Dashboard:  http://localhost:%s/ui\n", port)
	fmt.Printf("  API:        http://localhost:%s/api\n", port)
	fmt.Printf("  Data:       %s\n", dataDir)
	fmt.Printf("  ─────────────────────────────────\n\n")

	log.Printf("registry: listening on :%s", port)
	if err := http.ListenAndServe(":"+port, srv); err != nil {
		log.Fatalf("registry: %v", err)
	}
}
