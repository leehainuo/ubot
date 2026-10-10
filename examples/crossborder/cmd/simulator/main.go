package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"crossborder/internal/handler"
	"crossborder/internal/svc"
)

func main() {
	addr := os.Getenv("CROSSBORDER_ADDR")
	if addr == "" {
		addr = ":8091"
	}

	svcCtx := svc.NewServiceContext()
	mux := http.NewServeMux()
	handler.Init(mux, svcCtx)

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Crossborder simulator running on %s", addr)

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
