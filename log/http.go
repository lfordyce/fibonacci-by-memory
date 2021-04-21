package log

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"runtime/debug"
	"time"
)

func (l Logger) LoggerMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := l.Logger.With().Logger()

			mw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func(begin time.Time) {
				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.Error().
						Str("type", "error").
						Timestamp().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("log system error")
					http.Error(mw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				log.Info().
					Str("type", "access").
					Timestamp().
					Fields(map[string]interface{}{
						"remote_ip":   r.RemoteAddr,
						"url":         r.URL.Path,
						"proto":       r.Proto,
						"http_method": r.Method,
						"user_agent":  r.Header.Get("User-Agent"),
						"status":      mw.Status(),
						"latency_ms":  time.Since(begin).Milliseconds(),
						"bytes_in":    r.Header.Get("Content-Length"),
						"bytes_out":   mw.BytesWritten(),
					}).
					Msg("incoming_request")

			}(time.Now().UTC())

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
