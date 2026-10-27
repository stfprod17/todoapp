package users_transport_http

import (
	"net/http"

	core_logger "github.com/stfprod17/todoapp/internal/core/logger"
	core_http_request "github.com/stfprod17/todoapp/internal/core/transport/http/request"
	core_http_response "github.com/stfprod17/todoapp/internal/core/transport/http/response"
)

type LoginUserResponse struct {
	User  UserDTOResponse `json:"user"`
	Token string          `json:"access_token"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

func (h *UsersHTTPHandler) LoginUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request LoginUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate login request")
		return
	}

	user, err := h.usersService.LoginUser(ctx, request.Email, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login user")
		return
	}

	token, err := h.tokenManager.Generate(user)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to generate access token")
		return
	}

	response := LoginUserResponse{
		User:  userDTOFromDomain(user),
		Token: token,
	}

	responseHandler.JSONResponse(response, http.StatusOK)
}
