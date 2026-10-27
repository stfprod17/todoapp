package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/stfprod17/todoapp/internal/core/logger"
	core_http_response "github.com/stfprod17/todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			rw.Header().Set(requestIDHeader, requestID)
			next.ServeHTTP(rw, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			requestLogger := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)
			ctx := core_logger.ToContext(r.Context(), requestLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			rw := core_http_response.NewResponseWriter(w)
			before := time.Now()
			log.Debug(">>> incoming HTTP request", zap.String("http_method", r.Method), zap.Time("time", before.UTC()))
			next.ServeHTTP(rw, r)
			log.Debug("<<< doneHTTP request", zap.Int("status_code", rw.GetStatusCode()), zap.Duration("latency", time.Since(before)))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle HTTP request got unexpected panic")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
