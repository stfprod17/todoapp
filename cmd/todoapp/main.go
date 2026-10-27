package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_auth "github.com/stfprod17/todoapp/internal/core/auth"
	core_config "github.com/stfprod17/todoapp/internal/core/config"
	statistics_postgres_repository "github.com/stfprod17/todoapp/internal/core/features/statistics/repository/postgres"
	statistics_service "github.com/stfprod17/todoapp/internal/core/features/statistics/service"
	statistics_transport_http "github.com/stfprod17/todoapp/internal/core/features/statistics/transport/http"
	tasks_postgres_repository "github.com/stfprod17/todoapp/internal/core/features/tasks/repository/postgres"
	tasks_service "github.com/stfprod17/todoapp/internal/core/features/tasks/service"
	tasks_transport_http "github.com/stfprod17/todoapp/internal/core/features/tasks/transport/http"
	users_postgres_repository "github.com/stfprod17/todoapp/internal/core/features/users/repository/postgres"
	users_service "github.com/stfprod17/todoapp/internal/core/features/users/service"
	users_transport_http "github.com/stfprod17/todoapp/internal/core/features/users/transport/http"
	core_logger "github.com/stfprod17/todoapp/internal/core/logger"
	core_pgx_pool "github.com/stfprod17/todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/stfprod17/todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/stfprod17/todoapp/internal/core/transport/http/server"
	"go.uber.org/zap"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	fmt.Println("Hello todoapp")

	logger, err := core_logger.NewLogger(
		core_logger.NewConfigMust(),
	)
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug(
		"application time zone",
		zap.Any("zone", time.Local),
	)

	logger.Debug("initializing postgres connection pool")

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal(
			"failed to init postgres connection pool",
			zap.Error(err),
		)
	}
	defer pool.Close()

	logger.Debug(
		"initializing feature",
		zap.String("feature", "users"),
	)

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)

	tokenManager, err := core_auth.NewTokenManager(os.Getenv("JWT_SECRET"))
	if err != nil {
		logger.Fatal("failed to init token manager", zap.Error(err))
	}

	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(
		usersService,
		tokenManager,
	)

	logger.Debug(
		"initializing feature",
		zap.String("feature", "tasks"),
	)

	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(
		tasksService,
		tokenManager,
	)

	logger.Debug(
		"initializing feature",
		zap.String("feature", "statistics"),
	)

	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(
		pool,
	)
	statisticsService := statistics_service.NewStatisticsService(
		statisticsRepository,
	)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(
		statisticsService,
		tokenManager,
	)

	logger.Debug("initializing HTTP server")

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(
		core_http_server.ApiVersion1,
	)

	apiVersionRouter.RegisterRoutes(
		usersTransportHTTP.Routes()...,
	)

	apiVersionRouter.RegisterRoutes(
		tasksTransportHTTP.Routes()...,
	)

	apiVersionRouter.RegisterRoutes(
		statisticsTransportHTTP.Routes()...,
	)

	httpServer.RegisterAPIRouters(
		apiVersionRouter,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error(
			"HTTP server run error",
			zap.Error(err),
		)
	}
}
