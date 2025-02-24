package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)




func sendJson(w http.ResponseWriter, response Response, status int) {
	data, err := json.Marshal(response)
	if err != nil {
		slog.Error("error in JSON marshal", "error", err)
		return
	}
	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("error in write response", "error", err)
		return
	}
}

func NewHandler() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	r.Route("/api", func(r chi.Router) {
		r.Post("/short", handlePost)
		r.Get("/{code}", handleGet)
	})
	

	return r
}

func handlePost(w http.ResponseWriter, r *http.Request){

}

func handleGet(w http.ResponseWriter, r *http.Request) {
	// code := chi.URLParam(r, "code")
}


type PostBody struct {
	URL string `json:"url"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data any `json:"data,omitempty"`
}