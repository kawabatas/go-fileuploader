package handler

import (
	"net/http"
)

func (mux *ServerMux) routingPatterns(publicDir string) {
	mux.router.HandleFunc("GET /", mux.handleGetList)
	mux.router.HandleFunc("GET /files/{id}", mux.handleGetDetail)
	mux.router.HandleFunc("DELETE /files/{id}", mux.handleDelete)
	mux.router.HandleFunc("POST /files", mux.handlePost)
	mux.router.Handle("GET /images/", http.FileServer(http.Dir(publicDir)))
	mux.router.Handle("GET /debug/", http.FileServer(http.Dir(publicDir)))
}
