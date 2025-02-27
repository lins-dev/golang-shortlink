package api

import (
	"encoding/json"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func sendJson(w http.ResponseWriter, response Response, status int) {
	w.Header().Set("Content-Type", "application/json")
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

func NewHandler(db map[string]string) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	r.Route("/api", func(r chi.Router) {
		r.Post("/short", handlePost(db))
		r.Get("/{code}", handleGet(db))
		r.Get("/", handleIndex(db))
	})

	return r
}

func handlePost(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body PostBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJson(w, Response{Error: "invalid body"}, http.StatusUnprocessableEntity)
			return
		}
		if _, err := url.Parse(body.URL); err != nil {
			sendJson(w, Response{Error: "invalid URL"}, http.StatusBadRequest)
			return
		}
		code := getCode()
		db[code] = body.URL
		sendJson(w, Response{Data: code}, http.StatusCreated)
	}
}

func handleGet(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var code = chi.URLParam(r, "code")
		data, ok := db[code]
		if !ok {
			sendJson(w, Response{Error: "code not found"}, http.StatusNotFound)
			return
		}
		http.Redirect(w, r, data, http.StatusPermanentRedirect)
	}
}

func handleIndex(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendJson(w, Response{Data: db}, http.StatusOK)
	}
}

func getCode() string {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const limit = 8
	bytes := make([]byte, limit)

	for i := range limit {
		bytes[i] = characters[rand.IntN(len(characters))]
	}
	return string(bytes)
}

type PostBody struct {
	URL string `json:"url"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}
