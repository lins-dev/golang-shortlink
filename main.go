package main

import (
	"golang-shortlink/api"
	"log/slog"
	"net/http"
	"time"
)

func main() {

	if err := run(); err != nil {
		slog.Error("error in server", "error", err)
		return
	}
	slog.Info("server started")

}

func run() error {
	fakeDb := make(map[string]string)
	handler := api.NewHandler(fakeDb)
	server := http.Server{
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler,
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
