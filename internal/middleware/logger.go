package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(
	code int,
) {

	rw.statusCode = code

	rw.ResponseWriter.WriteHeader(code)
}

func Logger(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		fmt.Printf(
			"[%s] %s %s %d %v\n",
			time.Now().Format(
				"2006-01-02 15:04:05",
			),
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
		)
	})
}
