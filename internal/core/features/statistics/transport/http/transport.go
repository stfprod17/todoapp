package statistics_transport_http

import (
	"context"
	"net/http"

	core_auth "github.com/stfprod17/todoapp/internal/core/auth"
	core_http_middleware "github.com/stfprod17/todoapp/internal/core/transport/http/middleware"
	"time"

	"github.com/stfprod17/todoapp/internal/core/domain"
	core_http_server "github.com/stfprod17/todoapp/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
	tokenManager      *core_auth.TokenManager
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
	tokenManager *core_auth.TokenManager,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
		tokenManager:      tokenManager,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/statistics",
			Handler:    h.GetStatistics,
			Middleware: []core_http_middleware.Middleware{core_http_middleware.Auth(h.tokenManager)},
		},
	}
}
