package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Logger logs each request and recovers from panics.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				sw.status = http.StatusInternalServerError
				http.Error(sw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			log.Printf("%s %s %d %s", r.Method, r.URL.Path, sw.status, time.Since(start))
		}()
		next.ServeHTTP(sw, r)
	})
}
