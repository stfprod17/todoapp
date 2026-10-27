package users_transport_http

import (
	"context"
	"net/http"

	core_auth "github.com/stfprod17/todoapp/internal/core/auth"
	"github.com/stfprod17/todoapp/internal/core/domain"
	core_http_server "github.com/stfprod17/todoapp/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
	tokenManager *core_auth.TokenManager
}

type UsersService interface {
	CreateUser(ctx context.Context, user domain.User, password string) (domain.User, error)
	LoginUser(ctx context.Context, email string, password string) (domain.User, error)
	GetUsers(ctx context.Context, limit *int, offset *int) ([]domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	PatchUser(ctx context.Context, id int, patch domain.UserPatch) (domain.User, error)
}

func NewUsersHTTPHandler(usersService UsersService, tokenManager *core_auth.TokenManager) *UsersHTTPHandler {
	return &UsersHTTPHandler{usersService: usersService, tokenManager: tokenManager}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{Method: http.MethodPost, Path: "/users", Handler: h.CreateUser},
		{Method: http.MethodPost, Path: "/auth/login", Handler: h.LoginUser},
		{Method: http.MethodGet, Path: "/users", Handler: h.GetUsers},
		{Method: http.MethodGet, Path: "/users/{id}", Handler: h.GetUser},
		{Method: http.MethodDelete, Path: "/users/{id}", Handler: h.DeleteUser},
		{Method: http.MethodPatch, Path: "/users/{id}", Handler: h.PatchUser},
	}
}
