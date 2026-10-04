package main

import (
	"crossborder/internal/httpapi"
	"crossborder/internal/inventory"
	"crossborder/internal/order"
	"errors"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := os.Getenv("CROSSBORDER_ADDR")
	if addr == "" {
		addr = ":8091"
	}

	handler := httpapi.New(
		order.NewService(),
		inventory.NewService(time.Now),
	)

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Crossborder simulator running on %s", addr)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
