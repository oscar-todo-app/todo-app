package logs

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		start := time.Now()
		slog.Info(
			"url: %s, method: %s, addr: %s, duration: %s",
			r.RequestURI,
			r.Method,
			r.RemoteAddr,
			time.Since(start),
		)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
