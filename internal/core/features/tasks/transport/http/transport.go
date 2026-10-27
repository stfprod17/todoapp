package tasks_transport_http

import (
	"context"
	"net/http"

	core_auth "github.com/stfprod17/todoapp/internal/core/auth"
	"github.com/stfprod17/todoapp/internal/core/domain"
	core_http_middleware "github.com/stfprod17/todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/stfprod17/todoapp/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
	tokenManager *core_auth.TokenManager
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)

	GetTasks(
		ctx context.Context,
		userID *int,
		limit *int,
		offset *int,
	) ([]domain.Task, error)

	GetTask(
		ctx context.Context,
		id int,
	) (domain.Task, error)

	DeleteTask(
		ctx context.Context,
		id int,
	) error

	PatchTask(
		ctx context.Context,
		id int,
		patch domain.TaskPatch,
	) (domain.Task, error)
}

func NewTasksHTTPHandler(
	tasksService TasksService,
	tokenManager *core_auth.TokenManager,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
		tokenManager: tokenManager,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/tasks",
			Handler:    h.CreateTask,
			Middleware: []core_http_middleware.Middleware{core_http_middleware.Auth(h.tokenManager)},
		},
		{
			Method:     http.MethodGet,
			Path:       "/tasks",
			Handler:    h.GetTasks,
			Middleware: []core_http_middleware.Middleware{core_http_middleware.Auth(h.tokenManager)},
		},
		{
			Method:     http.MethodGet,
			Path:       "/tasks/{id}",
			Handler:    h.GetTask,
			Middleware: []core_http_middleware.Middleware{core_http_middleware.Auth(h.tokenManager)},
		},
		{
			Method:     http.MethodDelete,
			Path:       "/tasks/{id}",
			Handler:    h.DeleteTask,
			Middleware: []core_http_middleware.Middleware{core_http_middleware.Auth(h.tokenManager)},
		},
		{
			Method:     http.MethodPatch,
			Path:       "/tasks/{id}",
			Handler:    h.PatchTask,
			Middleware: []core_http_middleware.Middleware{core_http_middleware.Auth(h.tokenManager)},
		},
	}
}
