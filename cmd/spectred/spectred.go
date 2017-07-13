package main

import (
	"github.com/zeten30/spectre/http"
)

func main() {
	server := &http.Server{
		Addr:         ":8080",
		Proxied:      false,
		DocumentRoot: ".",

		PasteService:             &mockPasteService{},
		RequestPermitterProvider: loggingPermitter{},
	}
	server.Listen()
}
