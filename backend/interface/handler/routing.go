package handler

import (
	"net/http"

	"github.com/kawabatas/go-fileuploader/interface/handler/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func (mux *ServerMux) routingPatterns(publicDir string) {
	// handleFunc is a replacement for mux.HandleFunc
	// which enriches the handler's HTTP instrumentation with the pattern as the http.route.
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		// Configure the "http.route" for the HTTP instrumentation.
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))

		handler = middleware.UseAccessLogging(handler)
		mux.router.Handle(pattern, handler)
	}

	handleFunc("GET /", mux.handleGetList)
	handleFunc("GET /files/{id}", mux.handleGetDetail)
	handleFunc("DELETE /files/{id}", mux.handleDelete)
	handleFunc("POST /files", mux.handlePost)
	mux.router.Handle("GET /images/", http.FileServer(http.Dir(publicDir)))
	mux.router.Handle("GET /debug/", http.FileServer(http.Dir(publicDir)))
	http.Handle("/metrics", promhttp.Handler())
}
