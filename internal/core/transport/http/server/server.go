package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/stfprod17/todoapp/internal/core/logger"
	core_http_middleware "github.com/stfprod17/todoapp/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouters(
	routers ...*APIVersionRouter,
) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		handler := core_http_middleware.ChainMiddleware(
			router,
			router.middleware...,
		)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, handler),
		)
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	handler := core_http_middleware.ChainMiddleware(
		s.mux,
		s.middleware...,
	)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: handler,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn(
			"start HTTP server",
			zap.String("addr", s.config.Addr),
		)

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf(
				"listen and serve HTTP: %w",
				err,
			)
		}

	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf(
				"shutdown HTTP server: %w",
				err,
			)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
