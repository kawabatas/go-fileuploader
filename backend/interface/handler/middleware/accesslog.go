package middleware

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type loggingWriter struct {
	http.ResponseWriter
	code int
}

func newLoggingWriter(w http.ResponseWriter) *loggingWriter {
	return &loggingWriter{
		ResponseWriter: w,
		code:           http.StatusInternalServerError,
	}
}

func (lw *loggingWriter) WriteHeader(code int) {
	lw.code = code
	lw.ResponseWriter.WriteHeader(code)
}

var (
	meter      = otel.Meter("fileuploader-api-meter")
	requestCnt metric.Int64Counter
)

func init() {
	var err error
	requestCnt, err = meter.Int64Counter(
		"request",
		metric.WithDescription("Counts of Request for API"),
		metric.WithUnit("count"))
	if err != nil {
		panic(err)
	}
}

func UseAccessLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		rlw := newLoggingWriter(w)

		next.ServeHTTP(rlw, req)

		requestCnt.Add(req.Context(), 1)
		log.Printf("Doing really hard work requestCnt\n", 1)

		slog.InfoContext(req.Context(), "access log",
			slog.String("URI", req.RequestURI),
			slog.String("Method", req.Method),
			slog.Int("Status code", rlw.code),
			slog.Int64("Elapsed time (ms)", time.Since(start).Milliseconds()),
		)
	})
}
